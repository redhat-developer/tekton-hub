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

const Detail: React.FC = (props: any) => {
  const {taskId} = useParams();
  React.useEffect(() => {
    props.fetchTaskDescription(taskId);
    props.fetchTaskName(taskId);
    // eslint-disable-next-line
  }, []);

  const tempTask : any = [];
  if (props.TaskName != null) {
    tempTask.push(props.TaskName);
  }


  let taskDescription : string = '';
  if (props.TaskName != null) {
    taskDescription = (props.TaskName['description']);
  };
  const yamlData = '```'+props.TaskYaml+'```';
  return (
    <div>
      <Flex breakpointMods={[{modifier: 'row', breakpoint: 'lg'},
        {modifier: 'nowrap', breakpoint: 'lg'},
        {modifier: 'column', breakpoint: 'sm'}]}>
        <FlexItem>
          <Description
            Description = {props.TaskDescription}
            Yaml = {yamlData}
            userTaskDescription = {taskDescription} />
        </FlexItem>
      </Flex>
    </div>
  );
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

