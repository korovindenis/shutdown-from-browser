import React from 'react';
import HandleClick from './HandleClick';
import Countdown from 'react-countdown';
import string from 'prop-types';

class MyCountdown extends React.Component {
  render() {

    const renderer = ({ hours, minutes, seconds, completed }) => {
      if (completed) {
        HandleClick("shutdown", (new Date()).toISOString());
      } else {
        return <span>{hours}:{minutes}:{seconds}</span>;
      }
    };

    const dateNow = new Date()
    const time = dateNow.setTime(dateNow.getTime() + this.props.hours * 60 * 60 * 1000);
    
    return <Countdown date={time} renderer={renderer} />
  }
}

MyCountdown.propTypes = {
  hours: string
};

export default MyCountdown;