var index = require('../src/index.js');
var aws = require('aws-sdk-mock');
var fs = require('fs');

var chai = require('chai');
var should = chai.should();

var po = {
    srcBucket: 'ft-project-creator',
    dstBucket: 'ft-mio-projects',
    uuid: '83c42dbe-a9ec-4611-908b-2ffab9014042',
    name: 'some-fantastic-video-project'
};

describe('index', function(done) {
    it('should correctly munge project names', function(done) {

    });

    it('should create relevant files');
});
