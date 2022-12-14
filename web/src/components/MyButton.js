import React from 'react';
import Button from '@material-ui/core/Button';
import HandleClick from './HandleClick';
import {object, string} from 'prop-types';

class MyButton extends React.Component {
  render() {
    return <Button
      variant="contained"
      className={this.props.css}
      onClick={e => HandleClick(this.props.text, new Date(new Date().getTime()).toISOString().replace(/:[0-9]{2}\.[0-9]{3}Z/, ":00.000Z"))}
    >
      {this.props.text.toUpperCase()}
    </Button>;
  }
}

MyButton.propTypes = {
  css: string,
  text: string
};

export default MyButton;