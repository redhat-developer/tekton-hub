import React from 'react';
import {useParams} from 'react-router';
import {connect} from 'react-redux';
import BasicDetail from './BasicDetail';
import {fetchTaskName} from '../redux/Actions/TaskActionName';
import Loader from '../loader/loader';
import './basicdetail.css';

const Detail: React.FC = (props: any) => {
  const {taskId} = useParams();
  React.useEffect(() => {
    props.fetchTaskName(taskId);
    // eslint-disable-next-line
  }, []);
  if (props.TaskName && (props.TaskName.id).toString() === taskId) {
    return (
      <BasicDetail task={props.TaskName} />
    );
  }
  return (
    <Loader />
  );
};
const mapStateToProps = (state: any) => ({
  TaskName: state.TaskName.TaskName,
});
export default connect(mapStateToProps, {fetchTaskName})(Detail);

