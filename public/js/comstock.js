
function MyController($scope, $http) {
    $scope.response = ""
    $scope.submitData = function (user, resultVarName) {
		var mail = user.mail
		var password = user.password
		$http ({
			method:"GET",
			url: "/registerUser",
			params: {
				mail: mail,
				password: password
			}
		})
	    .success (function (data, status, headers, config) {
		$scope.response  = data.message
		$scope.isSucc = true
	    })
	    .error ( function (data, status, headers, config) {
		$scope.response = data.message
				$scope.isSucc = false
	    });
    }
    $scope.getIsSuccess = function() {
	return $scope.isSucc;
    }

}

