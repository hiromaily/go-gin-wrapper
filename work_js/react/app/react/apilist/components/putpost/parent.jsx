import React     from 'react'

import Put       from './put.jsx'
import Post      from './post.jsx'


export default class PutPostParent extends React.Component {
  render() {
    return (
      <div className='row'>
        <Put btn={this.props.btnPut} />
        <Post btn={this.props.btnPost} />
      </div>
    )
  }
}


/*
  var putBtn = document.getElementById("putBtn");
  var postBtn = document.getElementById("postBtn");
  putBtn.addEventListener("click", putUser, false);
  postBtn.addEventListener("click", postUser, false);


  function putUser(evt){
    var url;
    var sendData = new Object();
    var id = document.getElementById("putID").value;
    if (id != "" && !isNaN(id)){
        url = "/api/users/" + id;
    }else{
        //id error
        console.log("id is invalid.");
        swal("warning!", "user id is invalid.", "warning");
        //TODO:add has-error class to group
        return;
    }

    //create data
    if (document.getElementById("putFN").value != ""){
      sendData.firstName = document.getElementById("putFN").value;
    }
    if (document.getElementById("putLN").value != ""){
      sendData.lastName = document.getElementById("putLN").value;
    }
    if (document.getElementById("putEM").value != ""){
      sendData.email = document.getElementById("putEM").value;
    }
    if (document.getElementById("putPW").value != ""){
      sendData.password = document.getElementById("putPW").value;
    }

    console.log(sendData)

    //send
    hy.sendAjax(url, "put", "form", sendData)
  }

  function postUser(evt){
    var url = "/api/users";
    var sendData = new Object();

    var errFlg = 0;
    //create data
    if (document.getElementById("postFN").value != ""){
      sendData.firstName = document.getElementById("postFN").value;
    }else{
        errFlg=1;
    }
    if (document.getElementById("postLN").value != ""){
      sendData.lastName = document.getElementById("postLN").value;
    }else{
        errFlg=1;
    }
    if (document.getElementById("postEM").value != ""){
      sendData.email = document.getElementById("postEM").value;
    }else{
        errFlg=1;
    }
    if (document.getElementById("postPW").value != ""){
      sendData.password = document.getElementById("postPW").value;
    }else{
        errFlg=1;
    }

    //console.log(sendData)

    if (errFlg==1){
        swal("warning!", "blank filed is not allowed.", "warning");
        return;
    }

    //send
    hy.sendAjax(url, "post", "form", sendData)
  }

*/
