/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

var TableUsersController = function(users, $scope, $state, locationUtils) {

    $scope.users = users;

    $scope.editUser = function(id) {
        locationUtils.navigateToPath('/users/' + id);
    };

    $scope.create = function() {
        locationUtils.navigateToPath('/users/new');
    };

    $scope.register = function() {
        locationUtils.navigateToPath('/users/register');
    };

    $scope.refresh = function() {
        $state.reload(); // reloads all the resolves for the view
    };

    angular.element(document).ready(function () {
        $('#usersTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 25,
            "aaSorting": []
        });
    });

};

TableUsersController.$inject = ['users', '$scope', '$state', 'locationUtils'];
module.exports = TableUsersController;
