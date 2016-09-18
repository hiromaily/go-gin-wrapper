import React from 'react'
import Options from './options.jsx'

export default class Get extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)

    this.clickBtnEvt = this.clickBtnEvt.bind(this)
  }

  //Click get button
  clickBtnEvt(e) {
    console.log("[Get]:clickBtnEvt(e)")
    let idx = this.refs.slctIDs.selectedIndex
    let val = this.refs.slctIDs.options[idx].text

    //call event for post btn click
    this.props.btn.call(this, val)
  }

  render() {
    let options = this.props.ids.map(function (id) {
      return (
        <Options key={id} id={id} />
      )
    })

    return (
        <div className="col-md-6 col-sm-6 col-xs-12">
          <div className="panel panel-info">
            <div className="panel-heading">GET API</div>
            <div className="panel-body">
                <div className="form-group">
                  <label>Select user id</label>
                  <select id="slctIds" className="form-control" ref="slctIDs">
                    <option>All</option>
                    {options}
                  </select>
                </div>
                <button id="getBtn" type="button" onClick={this.clickBtnEvt} className="btn btn-info">Get User</button>
            </div>
          </div>
        </div>
    )
  }
}
