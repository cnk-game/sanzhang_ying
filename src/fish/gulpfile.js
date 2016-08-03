// 引入 gulp及组件
var gulp = require('gulp'), //基础库
    connect = require('gulp-connect');


/**
 * http 调试服务
 */
gulp.task('run', function() {
    // 开启服务器
    connect.server({
        root: './src',
	host: '192.168.1.133',
        port: 8088,
        livereload: true
    });
});
