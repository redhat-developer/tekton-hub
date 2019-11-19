import React from 'react';
import {connect} from 'react-redux';
import Task from '../task/Task';
import {fetchTaskSuccess} from '../redux/Actions/TaskAction';
import {Gallery} from '@patternfly/react-core';
import './index.css';

export interface MockData{
  Name : string,
  Description : string,
  Rating : number,
  Downloads : number,
  YAML : string
}

const TaskContainer: React.FC = (props: any) => {
  let tempArr : any = [];
  React.useEffect(() => {
    props.fetchTaskSuccess();
  }, []);

  if (props.TaskData != null) {
    tempArr = props.TaskData.map((task: any) =>{
      const taskData: MockData = {
        Name: task['name'],
        Description: task['description'],
        Rating: 0,
        Downloads: 0,
        YAML: task['yaml'],
      };
      return taskData;
    });
  }

  return (
    <div className="block">
      <Gallery gutter = "lg">
        {
          tempArr.map((task: any) => {
            return <Task key={task['name']} task = {task} />;
          })
        }
      </Gallery>
    </div>
  );

  // return (
  //   <div>
  //     {/* {
  //       tempArr.map((task: any) => {
  //         const taskData: MockData = {
  //           Name: task['Name'],
  //           Description: task['Description'],
  //           Rating: 0,
  //           Downloads: 0,
  //           YAML: task['YAML'],
  //         };
  //         // return <Task key={task['Name']} task={taskData} />;
  //       })
  //     } */}
  //     {
  //       tempArr.map((task: any) => {
  //         const taskData: MockData = {
  //           Name: task['Name'],
  //           Description: task['Description'],
  //           Rating: 0,
  //           Downloads: 0,
  //           YAML: task['YAML'],
  //         };
  //         return (<Gallery gutter = "md" key={task['Name']}><Task task = {taskData}/></Gallery>);
  //       })
  //     }
  //   </div>
  // );
};


const mapStateToProps = (state: any) => {
  return {
    TaskData: state.TaskData.TaskData,
  };
};

export default connect(mapStateToProps, {fetchTaskSuccess})(TaskContainer);


