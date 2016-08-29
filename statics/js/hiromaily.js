(function(){
  var hy = function(){};
  //public
  hy.debug = 1;

//-------------------------------------------------------------------------
//private
//-------------------------------------------------------------------------
  var updateUserList = function(users){
    var strHtml = "";
    for (var i = 0, len = users.length; i < len; i++) {
      //console.log(users[i].id, users[i].firstName);
      strHtml += "<tr><td>"+users[i].id+"</td><td>"+users[i].firstName+"</td><td>"+users[i].lastName+"</td>";
      strHtml += "<td>"+users[i].email+"</td><td>*****</td><td>"+users[i].update+"</td>";
    }
    var userListBody = document.getElementById("userListBody");
    userListBody.innerHTML = strHtml;
  };

  var setToken = function(token){
    var jwtCode = document.getElementById("jwtCode");
    jwtCode.value = token;
  };

  var getTokenHeader = function(){
    var jwtCode = document.getElementById("jwtCode");
    return jwtCode.value;
  };


//-------------------------------------------------------------------------
//public
//-------------------------------------------------------------------------
  //initialize
  hy.init = function(){};

  //ajax
  hy.sendAjax = function(url, method, content, sendData){
    var contentType = "application/x-www-form-urlencoded";
    if(content == "json"){
      contentType = "application/json";
      sendData = JSON.stringify(sendData);
    }

    var token = getTokenHeader();
    if (url != "/api/jwt" && token == ""){
      swal("error!", "token is required for sending ajax!", "error");
      return;
    } else if(url != "/api/jwt"){
      token = "Bearer " + token;
    }

    //for JSON
    $.ajax({
	  url: encodeURI(url),
	  type: method,
	  beforeSend: function (xhr) {
		//xhr.setRequestHeader('X-Custom-Header-Gin', '{{ .key }}');
		//xhr.setRequestHeader('{{ .header }}', '{{ .key }}');
		xhr.setRequestHeader(hiromaily_header, hiromaily_key);
        //'Authorization': 'Bearer ' + 'YOUR_CURRENT_TOKEN'
        if (token != ""){
		  xhr.setRequestHeader('Authorization', token);
        }
	  },
      //cache    : false,
      crossDomain: false,
      contentType: contentType,         //content of request body
      dataType   : 'json',              //data type from server
	  data:        sendData,
    })
    .done(function( data, textStatus, jqXHR ) {
      //console.log(JSON.stringify(data));
      console.log(data);
      if (method=="get" && data.code==0){
        updateUserList(data.users);
      }else if (method=="post" && data.token != null){
        console.log(data.token);
        setToken(data.token);
      }
      swal("success!", "user was updated!", "success");
    })
    .fail(function( jqXHR, textStatus, errorThrown ) {
      swal("error!", "validation error was occurred!", "error");
    });
  };

  window.hy = hy;
})();