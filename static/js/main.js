
function showNav(){
  snapper.open('left');
}


function dropExpand(e){
  var x = document.getElementById(e).style.display;
  console.log(x)
  if (x == 'block'){
    document.getElementById(e).style.display = 'none';
  } else{
    document.getElementById(e).style.display = 'block';
  }
}

function dropExpandInline(e){
  var x = document.getElementById(e).style.display;
  console.log(x)
  if (x == 'inline-block'){
    document.getElementById(e).style.display = 'none';
  } else{
    document.getElementById(e).style.display = 'inline-block';
  }
}


/***********************************************************************
angular
***********************************************************************/

;(function(){
function authInterceptor(API, auth, $location, $rootScope) {
  console.log($rootScope);

  return {
    // automatically attach Authorization header
    request: function(config) {
      var token = auth.getToken();
      if(token ) {
        //config.headers.Authorization = 'Bearer ' + token;
        config.headers['X-AUTH-TOKEN'] = token;
      }
      return config;
    },

    // If a token was sent back, save it
    response: function(res) {
      if(res.config.url.indexOf(API) === 0 && res.data.token) {
        auth.saveToken(res.data.token);
      }

      return res;
    }

  }
}

function authService($window) {
  var self = this;

    self.parseJwt = function(token) {
      var base64Url = token.split('.')[1];
      var base64 = base64Url.replace('-', '+').replace('_', '/');
      return JSON.parse($window.atob(base64));
    }

    self.saveToken = function(token) {
      $window.localStorage['jwtToken'] = token;
    }

    self.getToken = function() {
      return $window.localStorage['jwtToken'];
    }

    self.isAuthed = function() {
      var token = self.getToken();
      if(token) {
        var params = self.parseJwt(token);
        return Math.round(new Date().getTime() / 1000) <= params.exp;
      } else {
        return false;
      }
    }

    self.logout = function() {
      $window.localStorage.removeItem('jwtToken');
    }
  // Add JWT methods here
}

function userService($http, API, auth, $rootScope) {
  var self = this;
  self.user = null;
  self.roles = [];

  self.details = function() {
    return $http.get(API + '/me')
  }

  self.register = function(username, password) {
  return $http.post(API + '/auth/register', {
      username: username,
      password: password
    })
  }

  self.login = function(username, password, remember) {
  return $http.post(API + '/auth/login', {
      username: username,
      password: password,
      remember: remember,
    })
  };


  // add authentication methods here

}

var edna = angular.module('edna', ['ui.router', 'multiStepForm']);
edna.config(function($stateProvider, $urlRouterProvider) {
  //
  // For any unmatched url, redirect to /state1
  $urlRouterProvider.otherwise("/404");
  //
  // Now set up the states
  $stateProvider
    .state('404', {
      url: "/404",
      views: {
        "body": { templateUrl: "/partials/404.html" },
      },data:{
        roles: [],
        requireLogin: false,
      }
    })
    .state('login', {
      url: "/login",
      views: {
        "body": { templateUrl: "/partials/login.html" },
      },data:{
        roles: [],
        requireLogin: false,
      }
    })

    .state('root', {
      url: "/",
      views: {
        "content": {
          templateUrl: "/partials/dashboard.html",
          controller:function($rootScope, $sce){
            console.log("in root controller")

          },
        },
      },data:{
        roles: [],
        requireLogin: true,
      }
    })

    .state('root2', {
      url: "",
      views: {
        "content": { templateUrl: "/partials/dashboard.html" },
      },data:{
        roles: [],
        requireLogin: true,
      }
    })




        .state('teacher_list', {
          url: "/teacher/list",
          data:{
            roles: ['teacher'],
            requireLogin: true,
          },
          views: {
            "content": { templateUrl: "/partials/teacher/teacher_list.html" },
          },
        })

        .state('subjectTeacher', {
          abstract: true,
          url: "/teacher/subject/:id/:class",
          views: {
            "content": {
              templateUrl: "/partials/teacher/teacher_subject_root.html",
              controller:function($scope, $stateParams){
                console.log($stateParams)
                $scope.$stateParams = $stateParams;
                console.log("in teacher assessment area controller")

              },
             },
          },
          data:{
            roles: ['teacher'],
            requireLogin: true,
          }
        })
        .state('subjectTeacher.overview', {
          url: "/overview",
          views: {
            "teacher": { templateUrl: "/partials/teacher/teacher_subject.html" },
          },
        })
        .state('subjectTeacher.list', {
          url: "/list",
          views: {
            "teacher": { templateUrl: "/partials/teacher/teacher_subject_list.html" },
          },
        })

        .state('classTeacher', {
          abstract: true,
          url: "/teacher/class/:id",
          views: {
            "content": {
              templateUrl: "/partials/teacher/teacher_class_root.html",
              controller:function($scope, $stateParams){
                console.log($stateParams)
                $scope.$stateParams = $stateParams;
                console.log("in teacher assessment area controller")

              },
             },
          },
          data:{
            roles: ['teacher'],
            requireLogin: true,
          }
        })
        .state('classTeacher.overview', {
          url: "/overview",
          views: {
            "teacher": { templateUrl: "/partials/teacher/teacher_class.html" },
          },
        })
        .state('classTeacher.list', {
          url: "/list",
          views: {
            "teacher": { templateUrl: "/partials/teacher/teacher_class_list.html" },
          },
        })

    .state('staff', {
      views: {
        "content": {
          templateUrl: "/partials/staff/staff.html",
          controller:function($rootScope, $sce){
            console.log("in staff area controller")

          },
         },
      },data:{
        roles: ["admin"],
        requireLogin: true,
      }
    })
    .state('staff.new', {
      url: "/staff/new",
      views: {
        "staff": { templateUrl: "/partials/staff/staff_new.html" },
      },
    })
    .state('staff.list', {
      url: "/staff/list",

      views: {
        "staff": { templateUrl: "/partials/staff/staff_list.html" },
      },
    })
    .state('staff.edit', {
      url: "/staff/edit",

      views: {
        "staff": { templateUrl: "/partials/staff/staff_edit.html" },
      },
    })
    .state('staff.change_password', {
      url: "/staff/change_password",

      views: {
        "staff": { templateUrl: "/partials/staff/staff_change_password.html" },
      },
    })
    .state('staff_settings', {
      url: "/staff/settings",

      views: {
        "content": { templateUrl: "/partials/staff/staff_settings.html" },
      },data:{
        roles: [],
        requireLogin: true,
      }
    })
    .state('students', {
      views: {
        "content": { templateUrl: "/partials/students/students.html" },
      },data:{
        roles: ["admin"],
        requireLogin: true,
      }
    })
    .state('students.new', {
      url: "/students/new",
      views: {
        "students": { templateUrl: "/partials/students/students_new.html" },
      },
    })
    .state('students.list', {
      url: "/students/list",

      views: {
        "students": { templateUrl: "/partials/students/students_list.html" },
      },
    })
    .state('students.edit', {
      url: "/students/:id/edit",

      views: {
        "students": { templateUrl: "/partials/students/students_edit.html" },
      },data:{
        roles: ["admin", "teacher"],
        requireLogin: true,
      }
    })
    .state('students.viewresult', {
      url: "/students/:id/result",

      views: {
        "students": { templateUrl: "/partials/students/students_result.html" },
      },data:{
        roles: ["admin", "teacher"],
        requireLogin: true,
      }
    })
    .state('class', {
      url: "",
      views: {
        "content": { templateUrl: "/partials/class/class.html" },
      },data:{
        roles: ["admin"],
        requireLogin: true,
      }
    })
    .state('class.new', {
      url: "/class/new",
      views: {
        "class": { templateUrl: "/partials/class/class_new.html" },
      },
    })
    .state('class.list', {
      url: "/class/list",
      views: {
        "class": { templateUrl: "/partials/class/class_list.html" },
      },
    })
    .state('class.edit', {
      url: "/class/edit",
      views: {
        "class": { templateUrl: "/partials/class/class_edit.html" },
      },
    })
    .state('class.subject_new', {
      url: "/class/subjects/new",
      views: {
        "class": { templateUrl: "/partials/subject/subject_new.html" },
      },
    })
    .state('class.subject_list', {
      url: "/class/subjects/list",
      views: {
        "class": { templateUrl: "/partials/subject/subject_list.html" },
      },
    })
    .state('class.subject_edit', {
      url: "/class/subjects/edit",
      views: {
        "class": { templateUrl: "/partials/subject/subject_edit.html" },
      },
    })

    .state('settings_institution', {
      url: "/settings/institution",
      views: {
        "content": { templateUrl: "/partials/settings/institution_settings.html" },
      },data:{
        roles: ["admin"],
        requireLogin: true,
      }
    })
    .state('settings_session', {
      url: "/settings/session",
      views: {
        "content": { templateUrl: "/partials/settings/session_settings.html" },
      },data:{
        roles: ["admin"],
        requireLogin: true,
      }
    });
  });

  edna.factory('authInterceptor', authInterceptor)
  .service('user', userService)
  .service('auth', authService)
  .constant('API', '/api')
  .config(function($httpProvider) {
    $httpProvider.interceptors.push('authInterceptor');
  })
  .filter('assessmentsTotalFilter', function() {
    return function(studentA, overviewA ) {
      console.log(studentA);
      console.log(overviewA)
      var total = 0;
      for (i=0; i<studentA.length; i++){

        var upperLimit
        var percentage

        for (j=0; j<overviewA.length; j++){
          if (overviewA[i].name == studentA[i].name){
            upperLimit = parseInt(overviewA[i].upperlimit)
            percentage = parseInt(overviewA[i].percentage)
            break;
          }
        }


        //total = total + studentA[i].score
        console.log(studentA[i].score)
        console.log(upperLimit)

        xxx = (studentA[i].score /upperLimit)*(percentage)
        total = total + xxx

      }
      return total
    };
  })
  .filter('resultTotalFilter', function() {
    return function(studentA, overviewA ) {
      var total = 0;
      for (i=0; i<studentA.length; i++){

        var upperLimit
        var percentage

        for (j=0; j<overviewA.length; j++){
            upperLimit = parseInt(overviewA[i].upperlimit)
            percentage = parseInt(overviewA[i].percentage)
            break;

        }


        //total = total + studentA[i].score
        console.log(studentA[i].score)
        console.log(upperLimit)

        xxx = (studentA[i].score /upperLimit)*(percentage)
        total = total + xxx

      }
      return Math.round(total)
    };
  })
  .directive('restrict', function(user, $interpolate, $rootScope){
    return{
      restrict: 'A',
      priority: 100000,
      scope:true,
      link: function(scope, element, attr, linker){

        var findOne = function (haystack, arr) {
            return arr.some(function (v) {
                return haystack.indexOf(v) >= 0;
            });
        };

        //console.log(scope.x);

        var a = $interpolate(attr.access)(scope);
        console.log( a.trim() == "");
        if (a.trim() == ""){
          var attributes = []
        } else{
          var attributes = a.trim().split(" ");
        }

        if (user.roles.length == 0){
            user.details().then(function(res) {
              //console.log(res)
              $rootScope.user = res.data;
              user.user = res.data;
              user.roles = res.data.roles;

              var accessDenied = true;

              //console.log(res.data.roles);

              //console.log("vs");
              //console.log(attributes);
              if (findOne(res.data.roles, attributes)||attributes.length == 0){
                //console.log("Access denied in directive");
                accessDenied = false;
              }

              if (accessDenied){
                try {
                  element.children.remove();
                }catch(err){
                  console.log(err);
                }
                  //console.log(element)
                  //console.log("remove element");
                element.remove();
              }
            }, function (err){
              console.log("Error, user not authenticated");
              console.log(err);
            })
        }else{
          var accessDenied = true;

          //console.log(user.roles);
          //console.log("vs");
          //console.log(attributes);
          if (findOne(user.roles, attributes)||attributes.length == 0){
            console.log("Access denied in directive");
            accessDenied = false;
          }

          if (accessDenied){
            try {
              element.children.remove();
            }catch(err){
              console.log(err);
            }

            element.remove();
          }
        }
      },
    }
  })

  .run(function($rootScope, $state, auth, user, $sce){
    var dashboard = {
      nested:false,
      id:"Dashboard",
      name:"Dashboard",
      state: "root",
      roles: "",
      thumbnail: $sce.trustAsHtml('<i class="fa fa-home"></i>'),

    };

    var staff = {
      nested:true,
      id:"Staff",
      name:"Staff",
      state: "",
      roles:"admin",
      thumbnail: $sce.trustAsHtml('<i class="fa fa-user-secret "></i>'),
      children:[{
        id:"staff_new",
        name:"New",
        state:"staff.new",
        thumbnail:$sce.trustAsHtml('<i class="fa fa-plus"></i>'),
      },{
        id:"staff_list",
        name:"List",
        state:"staff.list",
        thumbnail:$sce.trustAsHtml('li'),
      }]
    };
    var classesnsubjects = {
      nested:true,
      id:"classesnsubjects",
      name:"Classes and Subjects",
      state: "",
      roles:"admin",
      thumbnail: $sce.trustAsHtml('<i class="fa fa-book"></i>'),
      children:[
      {
        id:"class_new",
        name:"New Class",
        state:"class.new",
        thumbnail:$sce.trustAsHtml('<i class="fa fa-plus"></i>'),
      },{
        id:"class_list",
        name:"Classes",
        state:"class.list",
        thumbnail:$sce.trustAsHtml('li'),
      },
      {
        id:"subject_new",
        name:"New Subjects",
        state:"class.subject_new",
        thumbnail:$sce.trustAsHtml('<i class="fa fa-plus"></i>'),
      },{
        id:"subject_list",
        name:"Subjects",
        state:"class.subject_list",
        thumbnail:$sce.trustAsHtml('li'),
      }]
    };

    var students = {
      nested:true,
      id:"Students",
      name:"Students",
      state: "",
      roles:"admin",
      thumbnail: $sce.trustAsHtml('<i class="fa fa-graduation-cap"></i>'),
      children:[{
        id:"students_new",
        name:"New",
        state:"students.new",
        thumbnail:$sce.trustAsHtml('<i class="fa fa-plus"></i>'),
      },{
        id:"students_list",
        name:"List",
        state:"students.list",
        thumbnail:$sce.trustAsHtml('li'),
      }]
    };


    var teacher = {
      nested:false,
      id:"AssignedClasses",
      name:"Classes",
      state:"teacher_list",
      roles:"teacher",
      thumbnail: $sce.trustAsHtml('<i class="fa fa-bookmark"></i>'),
    };

    var SchoolSettings = {
      nested:true,
      id:"Settings",
      name:"Settings",
      state: "",
      roles:"admin",
      thumbnail: $sce.trustAsHtml('<i class="fa fa-cogs"></i>'),
      children:[{
        id:"settings_institution",
        name:"Institution Settings",
        state:"settings_institution",
        thumbnail:$sce.trustAsHtml('<i class="fa fa-plus"></i>'),
      },{
        id:"settings_session",
        name:"Session Settings",
        state:"settings_session",
        thumbnail:$sce.trustAsHtml('li'),
      }]
    };

    $rootScope.addons = [dashboard, staff, classesnsubjects, students, teacher, SchoolSettings];

      console.log(auth.isAuthed())
      $rootScope.$on('$stateChangeStart', function(event, toState, toParams){
      var requireLogin = toState.data.requireLogin;
      var targetRoles = toState.data.roles;


      var findOne = function (haystack, arr) {
          return arr.some(function (v) {
              return haystack.indexOf(v) >= 0;
          });
      };

      console.log(user.user);

      if (requireLogin && auth.isAuthed()){
        if (user.user == null){
          console.log("user object is empty");
          user.details().then(function(res) {
            console.log(res)
            $rootScope.user = res.data;
            user.user = res.data;
            user.roles = res.data.roles;

          }, function (err){
            console.log("Error, user not authenticated")
            console.log(err)
            event.preventDefault();
            $state.go('login');
          })
        }else{
          console.log(user.roles);
          console.log(findOne(user.roles,targetRoles));

          console.log(targetRoles);
          if (findOne(user.roles,targetRoles) || targetRoles.length == 0 ){
            console.log("you can continue")

          }else{
            console.log("you are authenticated, but no required permissions");
            event.preventDefault();
            $state.go('404');
          }

      }
    }else if (requireLogin && !auth.isAuthed()){
      event.preventDefault();
      $state.go('login');
    }

    })

  })
  .controller('LoginCtrl', LoginCtrl)
  .controller('RootCtrl', RootCtrl)

  .controller('NewStaffCtrl', NewStaffCtrl)
  .controller('EditStaffCtrl', EditStaffCtrl)
  .controller('StaffChangePasswordCtrl', StaffChangePasswordCtrl)
  .controller('StaffListCtrl', StaffListCtrl)
  .controller('StaffSettingsCtrl', StaffSettingsCtrl)

  .controller('NewStudentCtrl', NewStudentCtrl)
  .controller('EditStudentCtrl', EditStudentCtrl)
  .controller('StudentListCtrl', StudentListCtrl)
  .controller('StudentResultCtrl', StudentResultCtrl)

  .controller('NewClassCtrl', NewClassCtrl)
  .controller('EditClassCtrl', EditClassCtrl)
  .controller('ClassListCtrl', ClassListCtrl)

  .controller('NewSubjectCtrl', NewSubjectCtrl)
  .controller('EditSubjectCtrl', EditSubjectCtrl)
  .controller('SubjectListCtrl', SubjectListCtrl)

  .controller('InstitutionSettingsCtrl', InstitutionSettingsCtrl)
  .controller('SessionSettingsCtrl', SessionSettingsCtrl)


  .controller('TeacherAssignedToCtrl', TeacherAssignedToCtrl)
  .controller('TeacherAssignedToSubjectCtrl', TeacherAssignedToSubjectCtrl)
  .controller('TeacherAssignedToSubjectOverviewCtrl', TeacherAssignedToSubjectOverviewCtrl)

  .controller('TeacherAssignedToClassCtrl', TeacherAssignedToClassCtrl)
  .controller('TeacherAssignedToClassOverviewCtrl', TeacherAssignedToClassOverviewCtrl);
})();
