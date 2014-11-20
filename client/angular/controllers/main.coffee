angular.module('redlion').controller 'mainController'
, ['$scope', '$q', '$rootScope', '$http'
, ($scope, $q, $rootScope, $http) ->


  currentOrder = null
  orderProcess = null

  ProcessOrder = (res, args) ->
    orderProcess = $q.defer()
    $rootScope.$broadcast 'order-processing', orderProcess.promise

    currentOrder.token = res.id
    currentOrder.args = args
    currentOrder.email = res.email
    $http.post('/api/v1/checkout', currentOrder)
    .success (data) ->
      orderProcess.resolve data
    .error () ->
      orderProcess.resolve {status: 'payment not completed'}

  GetPaymentDetails = (order) ->
    currentOrder = order
    StripeCheckout.open
      key:         redlion.stripePubKey
      amount:      $scope.grandTotal()
      currency:    'gbp'
      name:        'Red Lion After Hours'
      description: 'Secure Payment Form'
      billingAddress: true
      shippingAddress: $scope.order.post
      image:  '/img/logo.jpg'
      token:       ProcessOrder
    return

  $scope.orderError = false
  $scope.orderComplete = false

  $scope.order =
    quantity: 1
    post: false

  $scope.$on 'order-processing', (event, promise) ->
    $scope.orderProcessing = true
    $scope.hideOrderForm = true
    promise.then (result) ->
      $scope.orderProcessing = false
      if result.status is 'payment completed'
        $scope.orderComplete = true
        $scope.orderError = false
        $scope.hideOrderForm = true
        $scope.orderRef = result.id
      else
        $scope.hideOrderForm = false
        $scope.orderError = true
        $scope.orderComplete = false

  $scope.payWithStripe = () ->
    $scope.orderComplete = false
    $scope.orderError = false
    GetPaymentDetails($scope.order)

  $scope.postageCharge = () ->
    if $scope.order.post
      300
    else
      0

  $scope.processingFee = () ->
    $scope.order.quantity * 25 + 20

  $scope.grandTotal = () ->
    $scope.order.quantity * 1000 + $scope.postageCharge() + $scope.processingFee()

]
