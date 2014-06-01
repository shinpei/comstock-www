function MyController($scope, $http) {
    $scope.myData = {};
    $scope.myData.doClick = function(item, event) {
        var responsePromise = $http.get("/doc.html");
        responsePromise.success(function(data, status, headers, config) {
			console.log(data.title)
            $scope.myData.fromServer = data.title;
        });
        responsePromise.error(function(data, status, headers, config) {
			console.log(data)
        });
    }
}
