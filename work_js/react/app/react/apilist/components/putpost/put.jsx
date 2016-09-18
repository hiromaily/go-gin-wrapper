import React from 'react'

export default class Put extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
    this.state = {
        error: {
          all: '',
          id: '',
          firstName: '',
          lastName: '',
          email: '',
          pass: ''
        }
    }

    this.clickBtnEvt = this.clickBtnEvt.bind(this)
  }

  //Click put button
  //For validator and call passed method
  clickBtnEvt(e) {
    console.log("[Put]:clickBtnEvt(e)")
    //check values of form elements
    if (this.refs.putID.value == ""){
        //ID error
        this.setState({
          error: {all: "", id: "input id", firstName: "", lastName: "", email: "", pass: ""}
        })
        return
    }
    if (this.refs.putFN.value == "" && this.refs.putLN.value == "" &&
      this.refs.putEM.value == "" && this.refs.putPW.value == ""){
      this.setState({
        error: {all: "it requires change at least one", id:"", firstName: "", lastName: "", email: "", pass: ""}
      })
      return
    }
    //OK
    //reset error
    this.setState({
      error: {all: "", id:"", firstName: "", lastName: "", email: "", pass: ""}
    })

    //send data
    let sendData = [this.refs.putFN.value, this.refs.putLN.value, this.refs.putEM.value, this.refs.putPW.value]

    //call event for put btn click
    this.props.btn.call(this, sendData)
  }

  render() {
    return (
        <div className="col-md-6 col-sm-6 col-xs-12">
          <div className="panel panel-primary">
            <div className="panel-heading">PUT API</div>
            <div className="panel-body">
                <span className="form-err">{this.state.error.all}</span>
                <div className="form-group">
                  <label>User ID</label>
                  <input id="putID" className="form-control" type="text" ref="putID" />
                  <span className="form-err">{this.state.error.id}</span>
                </div>
                <div className="form-group">
                  <label>First Name</label>
                  <input id="putFN" className="form-control" type="text" ref="putFN" />
                  <span className="form-err">{this.state.error.firstName}</span>
                </div>
                <div className="form-group">
                  <label>Last Name</label>
                  <input id="putLN" className="form-control" type="text" ref="putLN" />
                  <span className="form-err">{this.state.error.lastName}</span>
                </div>
                <div className="form-group">
                  <label>E-mail</label>
                  <input id="putEM" className="form-control" type="text" autoComplete="off" ref="putEM" />
                  <span className="form-err">{this.state.error.email}</span>
                </div>
                <div className="form-group">
                  <label>Password</label>
                  <input id="putPW" className="form-control" type="password" autoComplete="off" ref="putPW" />
                  <span className="form-err">{this.state.error.pass}</span>
                </div>
                <button id="putBtn" type="button" onClick={this.clickBtnEvt} className="btn btn-info">Update</button>
            </div>
          </div>
        </div>
    )
  }
}





