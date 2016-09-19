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
