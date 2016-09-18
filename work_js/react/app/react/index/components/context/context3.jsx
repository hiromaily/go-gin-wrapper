import React from 'react'

class Provider extends React.Component {
  // Initialize some context
  getChildContext() {
    return {
      message: 'Stateless functional components can access context as well.'
    }
  }
  
  render() {
    return(
      <StatelessChild />
    )
  }
}

// Context-provider property
Provider.childContextTypes = {
  message: React.PropTypes.string
}

// Here we can access context inside our stateless child by passing it as the second argument to the function
const StatelessChild = (props, context) => {
  return <h2>{context.message}</h2>
}

// Context-user property
StatelessChild.contextTypes = {
  message: React.PropTypes.string
}
