module.exports = (grunt) ->
  grunt.initConfig
    clean:
      reset:
        src: ['bin', 'temp']

    cucumberjs:
      dev:
        src: 'test/e2e'
        options:
          format: 'pretty'

    shell:
      bowerInstall:
        command: 'bower install'
      buildUAT:
        command: 'goxc -pv=' + process.env.UAT_VER
        options:
          stdout: true
          stderr: true
      buildProd:
        command: 'goxc -pv=' + process.env.PROD_VER
        options:
          stdout: true
          stderr: true
      deployDev:
        command: 'ansible-playbook -i ' + process.env.ANSIBLE_HOSTS + ' deploy/dev.yml'
        options:
          stdout: true
          stderr: true
      deployUAT:
        command: 'ansible-playbook -i ' + process.env.ANSIBLE_HOSTS + ' deploy/uat.yml'
        options:
          stdout: true
          stderr: true
      deployProd:
        command: 'ansible-playbook -i ' + process.env.ANSIBLE_HOSTS + ' deploy/prod.yml'
        options:
          stdout: true
          stderr: true
    copy:
      UATServer:
        expand: true
        cwd: 'templates'
        src:['*.tmpl']
        dest: 'bin/releases/' + process.env.UAT_VER + '/linux_amd64/templates'
      UATPublic:
        expand: true
        cwd: 'public'
        src: ['**']
        dest: 'bin/releases/' + process.env.UAT_VER + '/linux_amd64/public'
      ProdServer:
        expand: true
        cwd: 'templates'
        src:['*.tmpl']
        dest: 'bin/releases/' + process.env.PROD_VER + '/linux_amd64/templates'
      ProdPublic:
        expand: true
        cwd: 'public'
        src: ['**']
        dest: 'bin/releases/' + process.env.PROD_VER + '/linux_amd64/public'
      fonts:
        expand: true
        cwd: 'bower_modules/bootstrap/'
        src: ['fonts/*']
        dest: 'public'
      images:
        expand: true
        cwd: 'client/'
        src: ['img/*']
        dest: 'public'

    less:
      styles:
        dest: 'public/css/styles.css'
        src: 'client/less/styles.less'

    coffee:
      scripts:
        expand: true
        cwd: 'client/angular'
        src: ['**/*.coffee']
        dest: 'temp/client/js/'
        ext: '.js'
        options:
          bare: true

    concat:
      scripts:
        src: [
          'bower_modules/underscore/underscore.js'
          'bower_modules/angular/angular.js'
          'bower_modules/angular-resource/angular-resource.js'
          'bower_modules/angular-route/angular-route.js'
          'temp/client/js/app.js'
          'temp/client/js/routes.js'
          'temp/client/js/services/*.js'
          'temp/client/js/directives/*.js'
          'temp/client/js/controllers/*.js'
          'temp/client/js/views.js'
          'temp/client/js/bootstrap.js'
        ]
        dest: 'public/js/script.js'

    ngTemplateCache:
      views:
        files:
          './temp/client/js/views.js': 'client/views/*.html'
        options:
          trim: 'client'
          module: 'redlion'

    uglify:
      scripts:
        files:
          'public/js/script.min.js': ['public/js/script.js']

    cssmin:
      styles:
        files:
          'public/css/styles.min.css': ['public/css/styles.css']
    watch:
      dev:
        files: ['server/**', 'client/**']
        tasks: ['build']
        options:
          livereload: true

  grunt.loadNpmTasks 'grunt-contrib-clean'
  grunt.loadNpmTasks 'grunt-contrib-coffee'
  grunt.loadNpmTasks 'grunt-contrib-concat'
  grunt.loadNpmTasks 'grunt-contrib-copy'
  grunt.loadNpmTasks 'grunt-contrib-less'
  grunt.loadNpmTasks 'grunt-contrib-cssmin'
  grunt.loadNpmTasks 'grunt-contrib-watch'
  grunt.loadNpmTasks 'grunt-contrib-uglify'
  grunt.loadNpmTasks 'grunt-gint'
  grunt.loadNpmTasks 'grunt-shell'

  grunt.registerTask 'build', [
    'clean'
    'coffee'
    'less:styles'
    'copy:fonts'
    'copy:images'
    'ngTemplateCache'
    'concat:scripts'
  ]

  grunt.registerTask 'buildUat', [
    'build'
    'uglify:scripts'
    'cssmin:styles'
  ]

  grunt.registerTask 'deployDev', [
    'shell:bowerInstall'
    'shell:deployDev'
  ]

  grunt.registerTask 'deployUat', [
    'buildUat'
    'shell:buildUAT'
    'copy:UATPublic'
    'copy:UATServer'
    'shell:deployUAT'
  ]

  grunt.registerTask 'deployProd', [
    'build'
    'uglify:scripts'
    'cssmin:styles'
    'clean:reset'
    'shell:buildProd'
    'copy:ProdPublic'
    'copy:ProdServer'
    'shell:deployProd'
  ]

  grunt.registerTask 'devServer', [
    'build'
    'watch'
  ]
