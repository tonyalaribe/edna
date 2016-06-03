
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

  function NewClassCtrl(API, $http, $scope, Notification) {
    console.log("new class");
    $scope.c = {};
    $scope.err = "";
    $scope.suc = "";
    $scope.c.teachers = [];
    $http.get(API + '/class').then(function(res) {
      console.log(res)
      $scope.classes = res.data.classes;

    }, function(err){
      console.log("Error")
      console.log(err)
    });

    $http.get(API + '/teachers').then(function(res) {
      console.log(res);
      $scope.teachers = res.data.users;

    }, function(err){
      console.log("Error")
      console.log(err)
    });


    function handleRequest(res) {
      console.log(res)
      $scope.c = {};
      //$scope.suc = "New Class Added Succesfully";
      Notification({message: 'New Class Added Succesfully', title: 'Class Management'});

    }

    function handleError(err){
      console.log("Error")
      console.log(err)
      Notification.error(err);
    }


    $scope.newclass = function(c){
      if(c.name){
        $scope.err = "";
        console.log(c);
        $scope.c = {};
        $http.post(API + '/class', c).then(handleRequest, handleError)
      }else {
        $scope.err = "Fill in Required Fields";
      }

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
      $rootScope.c = cl;
      $state.go("class.edit")
    }
  }


  function EditClassCtrl(API, $scope, $http, $state, $rootScope, Notification) {
    $scope.c = $rootScope.c;

    $http.get(API + '/teachers').then(function(res) {
      $scope.teachers = res.data.users;

    }, function(err){
      console.log("Error")
      console.log(err)
      Notification.error("Error Loading Teacher Info");
    });

    function handleRequest(res) {
      console.log(res)
      $scope.c = {};
      Notification({message: 'Class Data Successfuly Update!', title: 'Class Management'});
      $state.go("class.list")

    }

    function handleError(err){
      console.log("Error");
      console.log(err);
      Notification.error(err);
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
      Notification.error(err);
    });

  }



  /****************************************************
  Subject
  *****************************************************/


  function NewSubjectCtrl(API, $http, $scope, Notification) {
    $scope.subject = {};
    $scope.subject.teachers = [];
    $scope.err = "";
    $scope.suc = "";

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
      $scope.suc = "New Subject Added Succesfully";
      Notification({message: 'New Subject Added Succesfully', title: 'Subject Management'});

    }

    function handleError(err){
      console.log("Error")
      console.log(err)
      Notification.error(err);
    }


    $scope.newsubject = function(c){
      if(c.name){
        console.log(c);
        $scope.err = "";
        $http.post(API + '/subject', c).then(handleRequest, handleError)
      }else {
        $scope.err = "Fill in Required Fields";
      }

    }
  }


  function SubjectListCtrl(API, $scope, $http, $state, $rootScope) {
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
      $rootScope.subject = sub;
      $state.go("class.subject_edit")
    }
  }


  function EditSubjectCtrl(API, $scope, $http, $state, $rootScope, Notification) {
    $scope.subject = $rootScope.subject;


    $http.get(API + '/teachers').then(function(res) {
      console.log(res)
      $scope.teachers = res.data.users;

    }, function(err){
      console.log("Error")
      console.log(err)
      Notification.error(err);
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
      Notification({message: 'Subject Succesfully Updated', title: 'Subject Management'});
      $state.go("class.subject_list")

    }

    function handleError(err){
      console.log("Error")
      Notification.error(err);
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

  function NewStudentCtrl(API, $http, $scope, $location) {
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
      $scope.classes = res.data.classes;

    }, function(err){
      console.log("Error")
      console.log(err)
    });


    $scope.student = {};
    $scope.savestudent = function(student){
      $scope.submittedStudent = true;
      $http.post(API + '/student', student).then(
        function (res) {
          console.log(res)
          $scope.student = {};
          $scope.submittedStudent = false;
          $location.path("/students/list");
        }, function (err){
          console.log("Error")
          console.log(err)
          $scope.submittedStudent = false;
        }
      )
    }
  }

  function StudentListCtrl(API, $scope, $rootScope, $state, $http) {
    $http.get(API + '/students').then(
      function (res) {
        console.log(res)
        var students = res.data.students;
        $scope.students = students;
      }, function(err){
        console.log("Error")
        console.log(err)
      })

    $scope.edit = function(student){
      $rootScope.student = student;
      $state.go("students.edit", {id:student.id})
    }
  }

  function StudentResultCtrl(API, $scope, $rootScope, $state, $http, $stateParams) {
    $http.get(API + '/student?id='+$stateParams.id).then(
      function (res) {
        $scope.student = res.data.student;
      },function (err) {
      }
    )
    $http.get(API + '/student/result?id='+$stateParams.id).then(
      function (res) {
        var assessments = res.data;
        nassessment = []

        for (ii=0; ii<assessments.length; ii++){
          var objassess = assessments[ii].subjectinfo.assessments
          var stuassess = assessments[ii].assessments
          var returnArray = [];
          for (i=0; i<objassess.length; i++){

            for (j=0; j<stuassess.length; j++){
              if (objassess[i].name==stuassess[j].name){
                stuassess[j].upperlimit = objassess[i].upperlimit
                returnArray[i] = stuassess[j]
              }
            }


            if(!returnArray[i]){
              var  asStudent = {}
              asStudent.name = objassess[i].name
              asStudent.upperlimit = objassess[i].upperlimit
              asStudent.score = 0
              returnArray[i] = asStudent
            }

          }

          assessments[ii].assessments = returnArray
        }
        $scope.assessments = assessments;
      }, function(err){
        console.log("Error")
        console.log(err)
      })

  }


  function EditStudentCtrl(API, $scope, $http, $state, $rootScope, $stateParams) {
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
      $scope.classes = res.data.classes;

    }, function(err){
      console.log("Error")
      console.log(err)
    });

    console.log($stateParams.id);
    //$scope.student = $rootScope.student;
    //$scope.student.dateofbirth = new Date($scope.student.dateofbirth);
    //$scope.student.signupdate = new Date($scope.student.signupdate);

    $http.get(API + '/student?id='+$stateParams.id).then(function(res) {
      $scope.student = res.data.student;
      $scope.student.dateofbirth = new Date($scope.student.dateofbirth);
      $scope.student.signupdate = new Date($scope.student.signupdate);

    }, function(err){
      console.log("Error")
      console.log(err)
    });


    $scope.savestudent = function(student){
      console.log(student);
      $http.put(API + '/student', student).then(function(res) {
          console.log(res)
          $scope.student = {};
          $state.go("students.list")

        }, function (err){
            console.log("Error")
            console.log(err)
          }
        )
    }
  }



  /**********************************************************
  staff
  **********************************************************/

  function NewStaffCtrl(API, $http, $scope, Notification) {
    $scope.newstaff = {};
    $scope.res = "";
    function handleRequest(res) {
      console.log(res)
      $scope.newstaff = {};

      Notification({message: 'New Staff Data Added', title: 'Staff Management'});

    }


    function handleError(err){
      console.log("Error")
      console.log(err)
      Notification.error(err);
    }


    $scope.newstaffx = function(staff){
      if(staff.name && staff.type && staff.email && staff.phone && staff.password){
        $scope.res = "";
        $http.post(API + '/staff', staff).then(handleRequest, handleError);
      }else{
        $scope.res = "Please fill empty fields";
      console.log("error");
      }

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

    $scope.changepassword = function(staff){
      console.log("change password");
      console.log(staff);
      $rootScope.staff = staff;
      $state.go("staff.change_password");
    }
  }


  function EditStaffCtrl(API, $scope, $http, $state, $rootScope, Notification) {
    console.log("edit staff ctrl");
    $scope.staff = $rootScope.staff;
    function handleRequest(res) {
      console.log(res)
      Notification({message: 'Data Update Complete', title: 'Staff Management'});
      $scope.staff = {};
      $state.go("staff.list")

    }

    function handleError(err){
      console.log("Error")
      console.log(err)
      Notification.error(err);
    }


    $scope.editstaff = function(staff){
      console.log(staff);
      $http.put(API + '/staff', staff).then(handleRequest, handleError)
    }
  }

function StaffChangePasswordCtrl(API, $scope, $http, $state, $rootScope) {
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


  $scope.changepassword = function(staff){
    console.log(staff);

    $http.post(API + '/changepassword', staff).then(handleRequest, handleError)
  }
}

function StaffSettingsCtrl($scope, API, $http, $state, auth){

  $http.get(API + '/me').then(
    function(res){
      console.log(res.data);
      $scope.staff = res.data;
      $scope.f = res.data.image;
    }, function(err){console.log(err)})

    $scope.UpdateImage = function(file){
      var reader = new FileReader();
      reader.onload = function(u){
            $scope.$apply(function($scope) {
              $scope.f = u.target.result;
              $scope.staff.updateimage = u.target.result;
              //console.log(u.target.result);
            });
      };
      reader.readAsDataURL(file);

    };

    $scope.updateStaff = function(staff){
      $http.put(API + '/staff', staff).then(
        function(res){
          console.log(res)
          console.log("log out");
          auth.logout && auth.logout();
          $state.go("login")
        },
        function(err){
          console.log(err);
        });
    }
}

/********************************************************
/***********Settings**********************************/
/******************************************************/
function InstitutionSettingsCtrl(API, $scope, $http ){

  $http.get(API + '/sessions').then(function(res){
      console.log(res.data)
      $scope.sessions = res.data;
    },function(err){
      console.log(err)
    }
  );
  $http.get(API + '/school').then(function(res){
      console.log(res.data)
      $scope.institution = res.data;
    },function(err){
      console.log(err)
    }
  );


  $scope.updateInstitution = function(institution){

    $http.put(API + '/school', institution).then(
      function(res){
        console.log(res)
      }, function(err){
        console.log(err)
      }
    );
  }
}

function SessionSettingsCtrl(API, $scope, $http ){
  $http.get(API + '/sessions').then(function(res){
      console.log(res.data)
      $scope.sessions = res.data;
    },function(err){
      console.log(err)
    }
  );


  $scope.AddSession = function(session){
    session.id = session.start.getFullYear() + "/" + session.end.getFullYear()
    console.log(session)
    $http.post(API + '/session', session).then(
      function(res){
        console.log(res)
      }, function(err){
        console.log(err)
      }
    );
  }
}


/*******************************************************************
  //TEACHER TERRITORY
*******************************************************************/
  function TeacherAssignedToCtrl(API, $scope, $http ){
      $scope.AsSubjectTeacher = [];
      $scope.AsClassTeacher = [];
      $http.get(API + '/teacher/assignedto').then(function(res){
          console.log(res.data)
          $scope.AsSubjectTeacher = res.data.subjects;
          $scope.AsClassTeacher = res.data.classes;
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
            var assessmentLength = $scope.overview.assessments.length;

            students = res.data;
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

  function TeacherAssignedToSubjectOverviewCtrl(API, $http, $scope, $stateParams, Notification){
    $scope.overview = {};
    $scope.feedBack = "";
    $scope.assessment = {};
    $scope.indeterminate = "indeterminate";
    $scope.show = "hide";
    //console.log($scope.$parent.$stateParams)
    $http.get(API + '/subject?id='+encodeURI($stateParams.id)).then(function(res){
        console.log(res.data)
        $scope.overview = res.data.subject;
      },function(err){
        console.log(err)
      }
    );

    $scope.newAssessment = function(assessment){
      $scope.show = "show";
      $http.post(API + '/createassessment?id='+encodeURI($stateParams.id), assessment).then(function(res){
          console.log(res.data)
          $scope.show = "hide";
          $scope.feedBack = res.data;
            $scope.assessment = {};
            Notification({message: 'Success', title: 'Class Assesment'});
          $scope.overview.assessments.push(assessment);
        },function(err){
          console.log(err)
          $scope.show = "hide";
          Notification.error(err);
        }
      );
    }
  }


  function TeacherAssignedToClassCtrl(API, $scope, $http, $stateParams ){

    $scope.overview = {};
    console.log($stateParams)

    $scope.students = [];
    $http.get(API + '/studentsinclass?class='+encodeURI($stateParams.id)).then(function(res){
      console.log(res)
      $scope.students = res.data;
      },function(err){
        console.log(err)
      }
    );


  }

  function TeacherAssignedToClassOverviewCtrl(API, $http, $scope, $stateParams){
    $scope.overview = {};
    console.log($stateParams)

    $scope.students = [];
    $http.get(API + '/getclassdata?class='+encodeURI($stateParams.id)).then(function(res){
      console.log(res)

      $scope.overview.name = res.data.name
      $scope.overview.count = res.data.count
      console.log($scope.overview)
      //$scope.students = res.data;
      },function(err){
        console.log(err)
      }
    );
  }
