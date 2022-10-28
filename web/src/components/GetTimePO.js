import Component from 'react';
import string from 'prop-types';

class GetTimePO extends Component {
    constructor(){
        super()
        this.state = {
            loading: true,
            when: null, 
            error: null
        }
    }
    componentDidMount(){
        fetch('/api/v1/get-time-autopoweroff/', {
              headers: {
                'content-type': 'application/json',
                'accept': 'application/json',
              },
            })
            .then(res => res.json())
            .then(json => {
              this.setState({when: json.when, loading: false})
            })
            .catch(e => console.log(e))
      }

    render() {
      console.log('client check', this.state.isAuth)
      const { component: Component, ...props } = this.props;
      const {loading, when} = this.state;
      return (
        when
      )
    }
  }

  GetTimePO.propTypes = {
    component: string,
  };

  export default GetTimePO; 