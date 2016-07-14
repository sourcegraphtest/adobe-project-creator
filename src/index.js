var async = require("async");
var aws = require('aws-sdk');
var path = require("path");

var prefix = 'project';

exports.create = function(event, context, callback){
    var srcBucket = event.srcBucket;
    var dstBucket = event.dstBucket;

    var projectName = mutate(event.name);
    var projectUUID = event.uuid;

    var src = new aws.S3({params: {Bucket: srcBucket}, region: 'eu-west-1'});
    var dst = new aws.S3({params: {Bucket: dstBucket}, region: 'eu-west-1'});

    src.listObjects({Prefix: prefix}, function(err, data) {
        if(err) {
            console.log(err);
            callback();
        }

        if (data.Contents.length) {
            async.each(data.Contents, function(file, cb) {
                var baseFile = path.basename(file.Key);
                var suffix = path.extname(baseFile);

                if (baseFile != prefix) {
                    src.getObject({
                        Key: file.Key
                    }, function(err, data) {
                        if (err) {
                            console.log(err);
                            callback();
                        }
                        upload({bucket: dst,
                                projectName: projectName,
                                body: data.Body,
                                projectUUID: projectUUID,
                                suffix: suffix
                               },
                               cb);
                    });
                }

            }, function(err) {
                if (err) {
                    console.log(err);
                }
            });
        }
    });

    callback();

};

var upload = function(opts, callback) {
    var key = [opts.projectUUID, [opts.projectName, opts.suffix].join('')].join('/');

    opts.bucket.putObject({
        Key: key,
        Body: opts.body
    }, function(err, data) {
        if (err) {
            console.log(err);
        } else {
            callback();
        }
    });

};

var mutate = function(s) {
    return s.replace(/ /g, '_');
};


// var payload = {
//     srcBucket: 'ft-project-creator',
//     dstBucket: 'jspc-mio-s3-test',
//     uuid: 12345,
//     name: 'poo'
// };

// var x = exports.create(payload,
//                        {},
//                        function(err,out) {
//                            console.log(out);
//                        });
// console.log(x);
