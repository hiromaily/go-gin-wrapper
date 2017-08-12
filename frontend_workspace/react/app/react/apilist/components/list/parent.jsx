import React from 'react'
import Tables from './tables.jsx'

export default class List extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
  }


  render() {
    let tables = this.props.users.map(function (user) {
      return (
        <Tables key={user.id} user={user} />
      )
    })

    return (
      <div className='row'>
        <div className="col-lg-12 col-md-12 col-sm-12">
            <div className="table-responsive">
                <table className="table table-striped table-bordered table-hover">
                    <thead>
                    <tr>
                        <th>user_id</th>
                        <th>first_name</th>
                        <th>last_name</th>
                        <th>email</th>
                        <th>password</th>
                        <th>update_datetime</th>
                    </tr>
                    </thead>
                    <tbody id="userListBody">
                    {tables}
                    </tbody>
                </table>
            </div>
            <hr />
        </div>
      </div>
    )
  }
}
