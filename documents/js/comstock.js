
function MyController($scope, $http) {
	$scope.response = "HIHI"
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
				console.log("success")
				console.log($scope)
				$scope["response"]  = data.msg
				$scope["isSuccess"] = true
			})
			.error ( function (data, status, headers, config) {
				console.log("error")
				console.log($scope)
				$scope["response"] = data.msg
				$scope["isSuccess"] = false
			});
	}
}

