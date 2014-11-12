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
        command: 'ansible-playbook -i ' + process.env.ANSIBLE_HOSTS + ' deploy/dev/deploy.yml'
        options:
          stdout: true
          stderr: true
      deployUAT:
        command: 'ansible-playbook -i ' + process.env.ANSIBLE_HOSTS + ' deploy/uat/deploy.yml'
        options:
          stdout: true
          stderr: true
      deployProd:
        command: 'ansible-playbook -i ' + process.env.ANSIBLE_HOSTS + ' deploy/prod/deploy.yml'
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

    ngTemplateCache:
      views:
        files:
          './temp/client/js/views.js': 'client/views/*.html'
        options:
          trim: 'client'
          module: 'seasoned'

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
  grunt.loadNpmTasks 'grunt-contrib-copy'
  grunt.loadNpmTasks 'grunt-shell'
  grunt.loadNpmTasks 'grunt-contrib-less'
  grunt.loadNpmTasks 'grunt-contrib-cssmin'
  grunt.loadNpmTasks 'grunt-contrib-watch'

  grunt.registerTask 'build', [
    'clean'
    'less:styles'
    'copy:fonts'
    'copy:images'
  ]

  grunt.registerTask 'buildUat', [
    'build'
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
