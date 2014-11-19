angular.module('redlion').controller 'mainController'
, ['$scope'
, ($scope) ->

  $scope.model =
    message: "Pre orders available soon"

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
      amount:      $scope.grandTotal() * 100
      currency:    'gbp'
      name:        'Red Lion After Hours'
      description: 'Secure Payment Form'
      billingAddress: true
      shippingAddress: $scope.order.byPost
      image:  '/public/img/redlion.jpg'
      token:       ProcessOrder
    return

  $scope.orderError = false
  $scope.orderComplete = false

  $scope.order =
    quantity: 1
    byPost: false

  $scope.$on 'order-processing', (event, promise) ->
    $scope.orderProcessing = true
    promise.then (result) ->
      $scope.orderProcessing = false
      if result.status is 'payment completed'
        $scope.orderComplete = true
        $scope.orderError = false
      else
        $scope.orderError = true
        $scope.orderComplete = false

  $scope.payWithStripe = () ->
    $scope.orderComplete = false
    $scope.orderError = false
    GetPaymentDetails($scope.order)

  $scope.postageCharge = () ->
    if $scope.order.byPost
      3
    else
      0

  $scope.processingFee = () ->
    $scope.order.quantity * 0.25 + 0.2

  $scope.grandTotal = () ->
    $scope.order.quantity * 10 + $scope.postageCharge() + $scope.processingFee()

]
