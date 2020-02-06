import React from 'react';
import {connect} from 'react-redux';
import {Gallery,
  EmptyState,
  EmptyStateIcon,
  EmptyStateBody,
  EmptyStateVariant,
}
  from '@patternfly/react-core';
import Task from '../task/Task';
import {fetchTaskSuccess} from '../redux/Actions/TaskAction';
import {fetchTaskName} from '../redux/Actions/TaskActionName';
import './index.css';
import {CubesIcon} from '@patternfly/react-icons';
import Loader from '../loader/loader';

export interface TaskPropData{
  name : string,
  description : string,
  rating : number,
  downloads : number,
  yaml : string,
  tags : [],
}

const TaskContainer: React.FC = (props: any) => {
  let tempArr : any = [];
  React.useEffect(() => {
    props.fetchTaskSuccess();
  // eslint-disable-next-line
  }, []);

  if (props.TaskName != null) {
    for (let i = 0; i < props.TaskData.length; i++) {
      if (props.TaskName['id'] === props.TaskData[i]['id']) {
        tempArr.push(props.TaskData[i]);
      }
    }
  } else {
    if (props.TaskData != null) {
      tempArr = props.TaskData;

      if (tempArr.length === 0) {
        return (

          <div style = {{top: '50em',
            bottom: '50em', right: '50em', marginLeft: '45em'}}>
            <EmptyState variant={EmptyStateVariant.full}>
              <EmptyStateIcon icon={CubesIcon} />
              <EmptyStateBody>
          No match found.
              </EmptyStateBody>
            </EmptyState>
          </div>
        );
      }
    }
  }

  if (props.TaskData === undefined) {
    return (
      <div className="loader">
        <Loader />
      </div>
    );
  }

  return (
    <div className="block">
      <Gallery gutter="lg" style = {{marginRight: '-2.85em'}}>

        {
          tempArr.map((task: any) => <Task key={task.id} task={task} />)
        }


      </Gallery>
    </div>
  );
};


const mapStateToProps = (state: any) => ({
  TaskData: state.TaskData.TaskData,
  TaskName: state.TaskName.TaskName,

});
export default
connect(mapStateToProps, {fetchTaskSuccess, fetchTaskName})(TaskContainer);
