angular.module('redlion').config ['$routeProvider', '$locationProvider'
, ($routeProvider, $locationProvider) ->
  $locationProvider.html5Mode(true)
  .hashPrefix '!'
]
