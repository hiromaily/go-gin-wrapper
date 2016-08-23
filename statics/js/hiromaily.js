(function(){
  var hm = function(){};
  //public
  hm.debug = 1;

//-------------------------------------------------------------------------
//private
//-------------------------------------------------------------------------
  var updateUserList = function(users){
    var strHtml = "";
    for (var i = 0, len = users.length; i < len; i++) {
      console.log(users[i].id, users[i].firstName);
      strHtml += "<tr><td>"+users[i].id+"</td><td>"+users[i].firstName+"</td><td>"+users[i].lastName+"</td>";
      strHtml += "<td>"+users[i].email+"</td><td>*****</td><td>"+users[i].update+"</td>";
    }
    var userListBody = document.getElementById("userListBody");
    userListBody.innerHTML = strHtml;
  };

//-------------------------------------------------------------------------
//public
//-------------------------------------------------------------------------
  //initialize
  hm.init = function(){};

  //ajax
  hm.sendAjax = function(url, method, content, sendData){
    var contentType = "application/x-www-form-urlencoded";
    if(content == "json"){
      contentType = "application/json";
      sendData = JSON.stringify(sendData);
    }

    var rtnData;

    //for JSON
    $.ajax({
	  url: encodeURI(url),
	  type: method,
	  beforeSend: function (xhr) {
		//xhr.setRequestHeader('X-Custom-Header-Gin', '{{ .key }}');
		//xhr.setRequestHeader('{{ .header }}', '{{ .key }}');
		xhr.setRequestHeader(hiromaily_header, hiromaily_key);
	  },
      //cache    : false,
      crossDomain: false,
      contentType: contentType,         //content of request body
      dataType   : 'json',              //data type from server
	  data:        sendData,
    })
    .done(function( data, textStatus, jqXHR ) {
      //console.log("success");
      //console.log(JSON.stringify(data));
      //console.log(data);
      if (method=="get"){
        if(data.code==0){
          updateUserList(data.users);
        }
      } else if (method=="delete"){
        console.log("done delete");
      }
      swal("success!", "user was updated!", "success");
    })
    .fail(function( jqXHR, textStatus, errorThrown ) {
      console.log("error");
      swal("error!", "validation error was occurred!", "error");
    });
  };

  window.hm = hm;
})();