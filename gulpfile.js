var gulp = require('gulp');
var path = require('path');
var shell = require('gulp-shell');
var notifier   = require('node-notifier');
var child      = require('child_process');
var reload     = require('gulp-livereload');
var sync       = require('gulp-sync')(gulp).sync;
var util       = require('gulp-util');

var goPath = '/Users/user/src/go/src/github.com/maliceio/malice/**/*.go';

gulp.task('server:build', function() {
  var build = child.spawnSync('go', ['install']);
  if (build.stderr.length) {
    var lines = build.stderr.toString()
      .split('\n').filter(function(line) {
        return line.length;
      });
    for (var l in lines)
      util.log(util.colors.red(
        'Error (go install): ' + lines[l]
      ));
    notifier.notify({
      title: 'Error (go install)',
      message: lines
    });
  }
  return build;
});

gulp.task('server:watch', function() {

  /* Restart application server */
  gulp.watch([
    '.views/**/*.tmpl',
    'locales/*.json'
  ], ['server:spawn']);

  /* Rebuild and restart application server */
  gulp.watch([
    '*/**/*.go',
  ], sync([
    'server:build',
  ], 'server'));
});

gulp.task('build', [
  'server:build'
]);

/*
 * Start asset and server watchdogs and initialize livereload.
 */
gulp.task('watch', [
  'server:build'
], function() {
  reload.listen();
  return gulp.start([
    'server:watch',
  ]);
});

gulp.task('default', ['build']);
