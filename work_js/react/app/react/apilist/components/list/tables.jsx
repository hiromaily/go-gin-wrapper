import React from 'react'

export default class Options extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
  }

  render() {
    return (
        <tr>
            <th>{this.props.user.id}</th>
            <th>{this.props.user.firstName}</th>
            <th>{this.props.user.lastName}</th>
            <th>{this.props.user.email}</th>
            <th>*****</th>
            <th>{this.props.user.updated}</th>
        </tr>
    )
  }
}
