import React from 'react'

export default class JwtInput extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
  }

  render() {
    return (
      <div className="col-md-6 col-sm-6 col-xs-12">
        <div className="panel panel-primary">
          <div className="panel-heading">JWT API</div>
          <div className="panel-body">
            <form role="form">
              <div className="form-group">
                <label>E-mail</label>
                <input id="jwtEM" className="form-control" type="text" onChange={this.props.em} autoComplete="off" />
                <span className="form-err">{this.props.error.email}</span>
              </div>
              <div className="form-group">
                <label>Password</label>
                <input id="jwtPW" className="form-control" type="password" onChange={this.props.pw} autoComplete="off" />
                <span className="form-err">{this.props.error.pass}</span>
              </div>
              <button id="jwtBtn" type="button" onClick={this.props.btn} className="btn btn-info">Get JWT</button>
            </form>
          </div>
        </div>
      </div>
    )
  }
}





