import React from 'react';
import {useParams} from 'react-router';
import {connect} from 'react-redux';
import BasicDetail from './BasicDetail';
import {fetchTaskName} from '../redux/Actions/TaskActionName';
import {fetchTaskSuccess} from '../redux/Actions/TaskAction'
import Loader from '../loader/loader';
import './basicdetail.css';

const Detail: React.FC = (props: any) => {
  const {taskId} = useParams();
  React.useEffect(() => {
    props.fetchTaskSuccess();
    props.fetchTaskName(taskId);
    // eslint-disable-next-line
  }, []);

  if (props.TaskData) {
    for (let i = 0; i < props.TaskData.length; i++) {
      if (props.TaskData[i].id === Number(taskId)) {
        if (props.TaskName) {
          (props.TaskData[i]).data = props.TaskName.data
        }
        return (
          <BasicDetail task={props.TaskData[i]} />
        )
      }
    }
  }

  return (
    <Loader />
  );
};
const mapStateToProps = (state: any) => ({
  TaskName: state.TaskName.TaskName,
  TaskData: state.TaskData.TaskData,
});
export default connect(mapStateToProps, {fetchTaskName, fetchTaskSuccess})(Detail);

