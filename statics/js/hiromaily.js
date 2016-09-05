'use strict';

// hy object

var hy = new Object();
{
  (function () {
    //---------------------------------------------------------------------------
    //private
    //---------------------------------------------------------------------------
    var setToken = function setToken(token) {
      var jwtCode = document.getElementById('jwtCode');
      jwtCode.value = token;
    };

    var getTokenHeader = function getTokenHeader() {
      var jwtCode = document.getElementById('jwtCode');
      return jwtCode.value;
    };

    var updateUserList = function updateUserList(users) {
      console.info('updateUserList()');
      var strHtml = '';
      //for (let user of users) {
      users.forEach(function (user) {
        //console.log(users[i].id, users[i].firstName);
        strHtml += '<tr><td>' + user.id + '</td><td>' + user.firstName + '</td><td>' + user.lastName + '</td>\n<td>' + user.email + '</td><td>*****</td><td>' + user.update + '</td>';
      });
      var userListBody = document.getElementById('userListBody');
      userListBody.innerHTML = strHtml;
    };

    var square = function square(num) {
      return num * num;
    };

    //---------------------------------------------------------------------------
    //public
    //---------------------------------------------------------------------------
    //ajax
    hy.sendAjax = function (url, method, content, sendData) {
      var contentType = "application/x-www-form-urlencoded";
      if (content == "json") {
        contentType = "application/json";
        sendData = JSON.stringify(sendData);
      }

      var token = getTokenHeader();
      if (url != "/api/jwt" && token == "") {
        swal("error!", "token is required for sending ajax!", "error");
        return;
      } else if (url != "/api/jwt") {
        token = "Bearer " + token;
      }

      //for JSON
      $.ajax({
        url: encodeURI(url),
        type: method,
        beforeSend: function beforeSend(xhr) {
          //xhr.setRequestHeader('X-Custom-Header-Gin', '{{ .key }}')
          //xhr.setRequestHeader('{{ .header }}', '{{ .key }}')
          xhr.setRequestHeader(hiromaily_header, hiromaily_key);
          //'Authorization': 'Bearer ' + 'YOUR_CURRENT_TOKEN'
          if (token != "") {
            xhr.setRequestHeader('Authorization', token);
          }
        },
        //cache    : false,
        crossDomain: false,
        contentType: contentType, //content of request body
        dataType: 'json', //data type from server
        data: sendData
      }).done(function (data, textStatus, jqXHR) {
        //console.log(JSON.stringify(data));
        console.log(data);
        if (method == "get" && data.code == 0) {
          updateUserList(data.users);
        } else if (method == "post" && data.token != null) {
          console.log(data.token);
          setToken(data.token);
        }
        swal("success!", "user was updated!", "success");
      }).fail(function (jqXHR, textStatus, errorThrown) {
        swal("error!", "validation error was occurred!", "error");
      });
    };

    hy.abc = function () {
      console.log(square(5));
    };
  })();
}

// run
//hy.abc()
//# sourceMappingURL=hiromaily.js.map
