var async = require("async");
var aws = require('aws-sdk');
var fs = require('fs');
var http = require('http');
var mktemp = require('mktemp');
var path = require("path");
var strftime = require('strftime');
var strstr = require('string-to-stream');
var xml2js = require('xml2js');
var zlib = require('zlib');

var srcBucket = 'ft-project-creator';
var destBucket = 'jspc-mio-s3-test';

var s3 = new aws.S3();

exports.create = function(event, context, callback){
    var projectObj = {
        projectId: event.projectId,
        name: event.name,
        srcPath: event.srcPath,
        create: true
    };

    async.waterfall([
        async.apply(readProjectFile, projectObj),
        fromXml,
        mungeCreatedDate,
        mungeUpdatedDate,
        mungeLastPath,
        toXml,
        writeProjectFile
    ]);

    callback(null, {
        projectId: event.projectId,
        prproj: "https://s3-eu-west-1.amazonaws.com/jspc-mio-s3-test/" + event.projectId + "/" + event.name + ".prproj",
    });
};

var readProjectFile = function(po, cb) {
    var bucket,
        key;

    if (po.create) {
        bucket = srcBucket;
        key = "project/empty.prproj";
    } else {
        bucket = destBucket;
        key = id + "/" + name + ".prproj";
    }

    var gunzip = zlib.createGunzip();
    po.projXml = "";

    s3.getObject({Bucket: bucket, Key: key})
        .createReadStream()
        .pipe(gunzip)
        .on('data', function(chunk) {
            if (chunk) po.projXml += chunk.toString();
        })
        .on('close', function(err) {
            if (err) throw err;
            cb(null, po);
        });
};

var writeProjectFile = function(po, cb) {
    var gzip = zlib.createGzip();

    po.proj = "";
    strstr(po.projXml)
        .pipe(gzip)
        .on('data', function(chunk) {
        if (chunk) po.proj += chunk.toString();
        })
        .on('close', function(err) {
            if (err) throw err;

            s3.putObject({Bucket: destBucket,
                          Key: [po.projectId, [po.name, 'prproj'].join('.')].join('/'),
                          Body: po.proj
                         }, function(err, data) {
                             if (err) throw err;
                             else {
                                 console.log("Successfully uploaded data");
                                 cb(null, po);
                             }
                         });
        });

};

var fromXml = function(po, cb) {
    xml2js.parseString(po.projXml, function (err, result) {
        if (err) throw err;
        po.proj = result;

        cb(null, po);
    });
};

var toXml = function(po, cb) {
    var builder = new xml2js.Builder();
    po.projXml = builder.buildObject(po.proj);

    cb(null, po);
};

var mungeCreatedDate = function(po, cb) {
    po.proj.PremiereData.Project[1].Node[0].Properties[0]['MZ.BuildVersion.Created'] = "9.2.0x41 - " + strftime('%a %b %e %T %Y');

    cb(null, po);
};

var mungeUpdatedDate = function(po, cb) {
    po.proj.PremiereData.Project[1].Node[0].Properties[0]['MZ.BuildVersion.Modified'] = "9.2.0x41 - " + strftime('%a %b %e %T %Y');

    cb(null, po);
};

var mungeLastPath = function(po, cb) {
    po.proj.PremiereData.Project[1].Node[0].Properties[0]['project.settings.lastknowngoodprojectpath'] = path.join(po.srcPath, po.name, po.id + ".prproj");

    cb(null, po);
};
