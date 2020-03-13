import React from 'react';
import Description from '../description/Description';
import {connect} from 'react-redux';
import {
  Flex,
  FlexItem,
} from '@patternfly/react-core';
import {useParams} from 'react-router';
import {fetchTaskDescription} from '../redux/Actions/TaskActionDescription';
import {fetchTaskName} from '../redux/Actions/TaskActionName';
import store from '../redux/store';
const Detail: React.FC = (props: any) => {
  const {taskId} = useParams();
  React.useEffect(() => {
    // this dispatch is to reset previous description in redux store
    store.dispatch({type: 'FETCH_TASK_DESCRIPTION', payload: ''});
    props.fetchTaskDescription(taskId);
    props.fetchTaskName(taskId);
    // eslint-disable-next-line
  }, []);
  let taskDescription: string = '';
  let catalogTaskDescription: string = '';
  let yamlData: string = '';
  if (props.TaskDescription && (props.TaskName.id).toString() === taskId) {
    taskDescription = (props.TaskName['description']);
    catalogTaskDescription = props.TaskDescription;
    yamlData = '```' + props.TaskYaml + '```';
    return (
      <div>
        <Flex breakpointMods={[{modifier: 'row', breakpoint: 'lg'},
          {modifier: 'nowrap', breakpoint: 'lg'},
          {modifier: 'column', breakpoint: 'sm'}]}>
          <FlexItem>
            <Description
              Description={catalogTaskDescription}
              Yaml={yamlData}
              userTaskDescription={taskDescription} />
          </FlexItem>
        </Flex>
      </div>
    );
  } else {
    return (
      <div />
    );
  }
};

const mapStateToProps = (state: any) => {
  return {
    TaskDescription: state.TaskDescription.TaskDescription,
    TaskYaml: state.TaskYaml.TaskYaml,
    TaskName: state.TaskName.TaskName,
  };
};

export default
connect(mapStateToProps, {fetchTaskDescription, fetchTaskName})(Detail);

