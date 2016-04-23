
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

function RootCtrl(auth,user,$rootScope, $state) {
  var self = this;
  console.log("roor ctrl");

  document.getElementById('loader-wrapper').classList.add("loaded");


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

  self.dropExpand = dropExpand;
  self.dropExpandInline = dropExpandInline;

  self.logout = function() {
    console.log("log out");
    auth.logout && auth.logout();
    $state.go("login")
  }
  self.isAuthed = function() {
    return auth.isAuthed ? auth.isAuthed() : false
  }
}


/**************************************
LoginCtrl
***************************************/
function LoginCtrl(user, auth, $state, $rootScope, $scope) {
$scope.err = "";
  var self = this;
  self.remember = false;
  console.log("login ");



  function handleRequest2(res) {
    console.log(res)
    $rootScope.user = res.data;
  }

  function handleError2(err){
    console.log("Error")
    console.log(err)
  }




  function handleRequest(res) {
    console.log(res)
    var token = res.data.token ? res.data.token : null;
    if(token) {
      console.log('JWT:', token);
      $state.go("root");
      //location.reload();
      window.location.href="/";
      //user.details().then(handleRequest2, handleError2);
    }

  }

  function handleError(err){
    console.log("Error")
    $scope.err = "incorrect username or password";
    console.log(err)
  }

  self.login = function() {
    console.log("here");
    user.login(self.username, self.password, self.remember)
      .then(handleRequest, handleError)
  }

  self.logout = function() {
    auth.logout && auth.logout()
  }
  self.isAuthed = function() {
    return auth.isAuthed ? auth.isAuthed() : false
  }
}

/****
Class
****/

function NewClassCtrl(API, $http, $scope) {
  console.log("new class");
  $scope.c = {};
  $scope.c.teachers = [];

  $http.get(API + '/class').then(function(res) {
    console.log(res)
    $scope.classes = res.data.classes;

  }, function(err){
    console.log("Error")
    console.log(err)
  });

  $http.get(API + '/teachers').then(function(res) {
    console.log(res)
    $scope.teachers = res.data.users;

  }, function(err){
    console.log("Error")
    console.log(err)
  });


  function handleRequest(res) {
    console.log(res)
    $scope.c = {};

  }

  function handleError(err){
    console.log("Error")
    console.log(err)
  }


  $scope.newclass = function(c){
    console.log(c);
    $http.post(API + '/class', c).then(handleRequest, handleError)
  }
}


function ClassListCtrl(API, $scope, $http, $state, $rootScope) {
  console.log("class list ctrl");
  $scope.c = $rootScope.c;

  var teachers = [];

  $http.get(API + '/teachers').then(function(res) {
    console.log(res)
    teachers = res.data.users;

  }, function(err){
    console.log("Error getting teachers")
    console.log(err)
  });

  function handleRequest(res) {
    var x = res.data.classes;

    for (i = 0; i < x.length; i++){
      var t = x[i].teachers;
      var t2 = [];
      for (y = 0; y < t.length; y++){

        var elmpos = teachers.map(function(x){ console.log(x); return x.id;}).indexOf(t[y])
        t2.push(teachers[elmpos]);
      }

      x[i].teachers = t2;
    }
    $scope.classes = x;

  }

  function handleError(err){
    console.log("Error")
    console.log(err)
  }


  $http.get(API + '/class').then(handleRequest, handleError);

  $scope.edit = function(cl){
    console.log("edit");
    console.log(cl);
    $rootScope.c = cl;
    $state.go("class.edit")
  }
}


function EditClassCtrl(API, $scope, $http, $state, $rootScope) {
  console.log("edit class ctrl");
  $scope.c = $rootScope.c;

  $http.get(API + '/teachers').then(function(res) {
    console.log(res)
    $scope.teachers = res.data.users;

  }, function(err){
    console.log("Error")
    console.log(err)
  });

  function handleRequest(res) {
    console.log(res)
    $scope.c = {};
    $state.go("class.list")

  }

  function handleError(err){
    console.log("Error")
    console.log(err)
  }


  $scope.editclass = function(c){
    console.log(c);
    $http.put(API + '/class', c).then(handleRequest, handleError)
  }


  $http.get(API + '/class').then(function(res) {
    console.log(res)
    $scope.classes = res.data.classes;

  }, function(err){
    console.log("Error")
    console.log(err)
  });

}



/****************************************************
Subject
*****************************************************/

function NewSubjectCtrl(API, $http, $scope) {
  console.log("new subject");
  $scope.subject = {};
  $scope.subject.teachers = [];

  $http.get(API + '/teachers').then(function(res) {
    console.log(res)
    $scope.teachers = res.data.users;

  }, function(err){
    console.log("Error")
    console.log(err)
  });


  $http.get(API + '/class').then(function(res) {
    console.log(res)
    $scope.classes = res.data.classes;

  }, function(err){
    console.log("Error")
    console.log(err)
  });


  function handleRequest(res) {
    console.log(res)
    $scope.subject = {};

  }

  function handleError(err){
    console.log("Error")
    console.log(err)
  }


  $scope.newsubject = function(c){
    console.log(c);
    $http.post(API + '/subject', c).then(handleRequest, handleError)
  }
}


function SubjectListCtrl(API, $scope, $http, $state, $rootScope) {
  console.log("subject list ctrl");
  $scope.subject = $rootScope.subject;


  var teachers = [];

  $http.get(API + '/teachers').then(function(res) {
    console.log(res)
    teachers = res.data.users;

  }, function(err){
    console.log("Error getting teachers")
    console.log(err)
  });



  function handleRequest(res) {
    console.log(res)
    var x = res.data.subjects;

    for (i = 0; i < x.length; i++){
      var t = x[i].teachers;
      var t2 = [];
      for (y = 0; y < t.length; y++){

        var elmpos = teachers.map(function(x){ console.log(x); return x.id;}).indexOf(t[y])
        t2.push(teachers[elmpos]);
      }

      x[i].teachers = t2;
    }

    $scope.subjects = x;
  }

  function handleError(err){
    console.log("Error")
    console.log(err)
  }


  $http.get(API + '/subjects').then(handleRequest, handleError);

  $scope.edit = function(sub){
    console.log("edit");
    console.log(sub);
    $rootScope.subject = sub;
    $state.go("class.subject_edit")
  }
}


function EditSubjectCtrl(API, $scope, $http, $state, $rootScope) {
  console.log("edit subject ctrl");
  $scope.subject = $rootScope.subject;


  $http.get(API + '/teachers').then(function(res) {
    console.log(res)
    $scope.teachers = res.data.users;

  }, function(err){
    console.log("Error")
    console.log(err)
  });

  $http.get(API + '/class').then(function(res) {
    console.log(res)
    $scope.classes = res.data.classes;

  }, function(err){
    console.log("Error")
    console.log(err)
  });


  function handleRequest(res) {
    console.log(res)
    $scope.subject = {};
    $state.go("class.subject_list")

  }

  function handleError(err){
    console.log("Error")
    console.log(err)
  }


  $scope.editsubject = function(subject){
    console.log(subject);
    $http.put(API + '/subject', subject).then(handleRequest, handleError)
  }
}


/**********************************************************
staff
**********************************************************/

function NewStudentCtrl(API, $http, $scope) {
  console.log("new staff");

  $scope.steps = [
    {
        templateUrl: '/partials/students/students_new_official.html',
        title: 'Official Details',
        hasForm: true,
    },
    {
        templateUrl: '/partials/students/students_new_personal.html',
        title: 'Personal Details',
        hasForm: true,
    },
    {
        templateUrl: '/partials/students/students_new_contacts.html',
        title: 'Contact Details',
        hasForm: true,
    },
    {
        templateUrl: '/partials/students/students_new_guardians.html',
        title: 'Guardian Details',
        hasForm: true,
    },
    {
        templateUrl: '/partials/students/students_new_previousqualification.html',
        title: 'Previous Qualification Details',
        hasForm: true,

    }
];

$http.get(API + '/class').then(function(res) {
  console.log(res)
  $scope.classes = res.data.classes;

}, function(err){
  console.log("Error")
  console.log(err)
});


  $scope.student = {};
  function handleRequest(res) {
    console.log(res)
    $scope.student = {};
    $scope.submittedStudent = false;
  }

  function handleError(err){
    console.log("Error")
    console.log(err)
    $scope.submittedStudent = false;
  }


  $scope.newstudent = function(student){
    $scope.submittedStudent = true;
    $http.post(API + '/student', student).then(handleRequest, handleError)
  }
}

function StudentListCtrl(API, $scope, $rootScope, $state, $http) {
  console.log("students list ctrl");
  function handleRequest(res) {
    console.log(res)

    var students = res.data.students;
    $scope.students = students;


  }

  function handleError(err){
    console.log("Error")
    console.log(err)
  }


  $http.get(API + '/student').then(handleRequest, handleError)

  $scope.edit = function(student){
    console.log("edit");
    console.log(student);
    $rootScope.student = student;
    $state.go("students.edit")
  }
}


function EditStudentCtrl(API, $scope, $http, $state, $rootScope) {
  console.log("edit staff ctrl");


    $scope.steps = [
      {
          templateUrl: '/partials/students/students_new_official.html',
          title: 'Official Details',
          hasForm: true,
      },
      {
          templateUrl: '/partials/students/students_new_personal.html',
          title: 'Personal Details',
          hasForm: true,
      },
      {
          templateUrl: '/partials/students/students_new_contacts.html',
          title: 'Contact Details',
          hasForm: true,
      },
      {
          templateUrl: '/partials/students/students_new_guardians.html',
          title: 'Guardian Details',
          hasForm: true,
      },
      {
          templateUrl: '/partials/students/students_new_previousqualification.html',
          title: 'Previous Qualification Details',
          hasForm: true,

      }
  ];

  $http.get(API + '/class').then(function(res) {
    console.log(res)
    $scope.classes = res.data.classes;

  }, function(err){
    console.log("Error")
    console.log(err)
  });



  $scope.student = $rootScope.student;

  $scope.student.dateofbirth = new Date($scope.student.dateofbirth);
  $scope.student.signupdate = new Date($scope.student.signupdate);
  function handleRequest(res) {
    console.log(res)
    $scope.student = {};
    $state.go("students.list")

  }

  function handleError(err){
    console.log("Error")
    console.log(err)
  }


  $scope.newstudent = function(student){
    console.log(student);
    $http.put(API + '/student', student).then(handleRequest, handleError)
  }
}


/**********************************************************
staff
**********************************************************/

function NewStaffCtrl(API, $http, $scope) {
  console.log("new staff");
  $scope.newstaff = {};
  function handleRequest(res) {
    console.log(res)
    $scope.newstaff = {};

  }


  function handleError(err){
    console.log("Error")
    console.log(err)
  }


  $scope.newstaffx = function(staff){
    $http.post(API + '/staff', staff).then(handleRequest, handleError)
  }
}

function StaffListCtrl(API, $scope, $rootScope, $state, $http) {
  console.log("staff list ctrl");
  function handleRequest(res) {
    console.log(res)
    $scope.staff = res.data.users;

  }

  function handleError(err){
    console.log("Error")
    console.log(err)
  }


  $http.get(API + '/staff').then(handleRequest, handleError)

  $scope.edit = function(staff){
    console.log("edit");
    console.log(staff);
    $rootScope.staff = staff;
    $state.go("staff.edit")
  }
}


function EditStaffCtrl(API, $scope, $http, $state, $rootScope) {
  console.log("edit staff ctrl");
  $scope.staff = $rootScope.staff;
  function handleRequest(res) {
    console.log(res)
    $scope.staff = {};
    $state.go("staff.list")

  }

  function handleError(err){
    console.log("Error")
    console.log(err)
  }


  $scope.editstaff = function(staff){
    console.log(staff);
    $http.put(API + '/staff', staff).then(handleRequest, handleError)
  }
}


//TEACHER TTERRITORY
function TeacherAssignedToCtrl(API, $scope, $http ){

    $scope.AsSubjectTeacher = [];
    $scope.AsClassTeacher = [];
    $http.get(API + '/teacher/assignedto').then(function(res){
        console.log(res.data)
        $scope.AsSubjectTeacher = res.data.subjects;
      },function(err){
        console.log(err)
      }
    );


}

function TeacherAssignedToSubjectCtrl(API, $scope, $http, $stateParams ){

  $scope.overview = {};

  $http.get(API + '/subject?id='+encodeURI($stateParams.id)).then(function(res){
      console.log(res.data)
      $scope.overview = res.data.subject;
      $scope.overview.class = $stateParams.class;
      $scope.assessments = res.data.subject.assessments;
      $scope.students = [];
      $http.get(API + '/studentsinclass?class='+encodeURI($stateParams.class)).then(function(res){
      //$http.get(API + '/studentsinclass?class=JSS+1+B').then(function(res){
          console.log(res.data)
          var assessmentLength = $scope.overview.assessments.length;

          students = res.data;
          /*for (j=0; j<students.length; j++){
            for (i=0; i<$scope.overview.assessments.length; i++){

              //console.log(students[j])
              //console.log(students[j].assessments[i])

              if(!students[j].assessments[i]){
                students[j].assessments[i] = {}
                students[j].assessments[i].name = $scope.overview.assessments[i].name
                students[j].assessments[i].score = 0
              }
            }

          }
          */


          var overviewA = $scope.overview.assessments;

          for (jj=0; jj<students.length; jj++){
            var studentA = students[jj].assessments;
            console.log(studentA)
            var returnStudentA;
            var returnArray = [];
            for (i=0; i<overviewA.length; i++){

              for (j=0; j<studentA.length; j++){
                if (overviewA[i].name==studentA[j].name){
                  studentA[j].upperlimit = overviewA[i].upperlimit
                  returnArray[i] = studentA[j]
                }
              }


              if(!returnArray[i]){
                var  asStudent = {}
                asStudent.name = overviewA[i].name
                asStudent.upperlimit = overviewA[i].upperlimit
                asStudent.score = 0
                returnArray[i] = asStudent
              }

            }

            console.log(returnArray);

            students[jj].assessments = returnArray;
        }

        $scope.students = students
        },function(err){
          console.log(err)
        }
      );

    },function(err){
      console.log(err)
    }
  );




  $scope.updateAssessment = function(s, a){

    var assessment = {}
    //console.log($scope.overview)
    console.log(a)
    assessment.studentid = s.studentid
    assessment.name = s.name
    assessment.subject = $scope.overview.name
    assessment.class = $scope.overview.class
    assessment.assessmentname = a.name
    console.log(a.score)
    assessment.score = parseInt(a.score)

    console.log(assessment)

    $http.post(API + '/addstudentassessment', assessment).then(function(res){
        console.log(res.data)
        //$scope.asessments.push(assessment);
      },function(err){
        console.log(err)
      }
    );
  }

}

function TeacherAssignedToSubjectOverviewCtrl(API, $http, $scope, $stateParams){
  $scope.overview = {};
console.log($scope.$parent.$stateParams)
  $http.get(API + '/subject?id='+encodeURI($stateParams.id)).then(function(res){
      console.log(res.data)
      $scope.overview = res.data.subject;
    },function(err){
      console.log(err)
    }
  );

  $scope.newAssessment = function(assessment){
    $http.post(API + '/createassessment?id='+encodeURI($stateParams.id), assessment).then(function(res){
        console.log(res.data)
        //$scope.asessments.push(assessment);
      },function(err){
        console.log(err)
      }
    );
  }
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

        .state('teacher', {
          abstract: true,
          url: "/teacher/subject/:id/:class",
          views: {
            "content": {
              templateUrl: "/partials/teacher/teacher.html",
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
        .state('teacher.subject_overview', {
          url: "/overview",
          views: {
            "teacher": { templateUrl: "/partials/teacher/teacher_subject.html" },
          },
        })
        .state('teacher.subject_list', {
          url: "/list",
          views: {
            "teacher": { templateUrl: "/partials/teacher/teacher_subject_list.html" },
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
      url: "/students/edit",

      views: {
        "students": { templateUrl: "/partials/students/students_edit.html" },
      },
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
            upperLimit = overviewA[i].upperlimit
            percentage = overviewA[i].percentage
            break;
          }
        }


        //total = total + studentA[i].score

        xxx = (studentA[i].score /upperLimit)*(percentage)
        total = total + xxx

      }
      return total
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
      thumbnail: $sce.trustAsHtml('<i class="fa fa-plus"></i>'),
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
      thumbnail: $sce.trustAsHtml('<i class="fa fa-group"></i>'),
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
      thumbnail: $sce.trustAsHtml('<i class="fa fa-group"></i>'),
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
      thumbnail: $sce.trustAsHtml('<i class="fa fa-group"></i>'),
    };

    $rootScope.addons = [dashboard, staff, classesnsubjects, students, teacher];

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
  .controller('StaffListCtrl', StaffListCtrl)

  .controller('NewStudentCtrl', NewStudentCtrl)
  .controller('EditStudentCtrl', EditStudentCtrl)
  .controller('StudentListCtrl', StudentListCtrl)

  .controller('NewClassCtrl', NewClassCtrl)
  .controller('EditClassCtrl', EditClassCtrl)
  .controller('ClassListCtrl', ClassListCtrl)

  .controller('NewSubjectCtrl', NewSubjectCtrl)
  .controller('EditSubjectCtrl', EditSubjectCtrl)
  .controller('SubjectListCtrl', SubjectListCtrl)

  .controller('TeacherAssignedToCtrl', TeacherAssignedToCtrl)
  .controller('TeacherAssignedToSubjectCtrl', TeacherAssignedToSubjectCtrl)
  .controller('TeacherAssignedToSubjectOverviewCtrl', TeacherAssignedToSubjectOverviewCtrl);
})();
