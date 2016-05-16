var aws = require('aws-sdk');
var mktemp = require('mktemp');
var http = require('http');
var fs = require('fs');
var path = require("path");
var async = require("async");
var zlib = require('zlib');
var xml2js = require('xml2js');
var strftime = require('strftime');
var strstr = require('string-to-stream')

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
    if (po.create) {
        var bucket = srcBucket;
        var key = "project/empty.prproj";
    } else {
        var bucket = destBucket;
        var key = id + "/" + name + ".prproj";
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
    var w = fs.createWriteStream('./out.prproj');
    var gzip = zlib.createGzip();

    strstr(po.projXml)
        .pipe(gzip)
        .pipe(w);

    cb(null, po);
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
}

exports.create({projectId: 12345,
                name: 'my_project',
                srcPath: '/Users/jspc/tmp'},
               {},
               function(err, data) {
                   console.log(data);
               });
