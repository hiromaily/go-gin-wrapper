import React from 'react'

export default class JwtCode extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
  }

  render() {
    return (
        <div className="col-md-6 col-sm-6 col-xs-12">
          <div className="panel panel-info">
            <div className="panel-heading">JWT CODE</div>
            <div className="panel-body">
                <div className="form-group">
                  <label htmlFor="disabledInput">jwt code</label>
                  <input id="jwtCode" className="form-control" id="disabledInput" value={this.props.code} type="text" placeholder="" disabled="" />
                </div>
            </div>
          </div>
        </div>
    )
  }
}
