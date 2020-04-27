import {FETCH_TASK_DESCRIPTION} from '../Actions/TaskActionType';
import {FETCH_TASK_YAML} from '../Actions/TaskActionType';

// eslint-disable-next-line require-jsdoc
export function fetchTaskDescription(name: any, version: any) {

  return function (dispatch: any) {

    fetch(`https://raw.githubusercontent.com/Pipelines-Marketplace/catalog/master/official/tasks/${name}/v${version}/README.md`)
      .then((response) => response.text())
      .then((TaskDescription) => dispatch({
        type: FETCH_TASK_DESCRIPTION,
        payload: TaskDescription
      }))

    fetch(`https://raw.githubusercontent.com/Pipelines-Marketplace/catalog/master/official/tasks/${name}/v${version}/${name}.yaml`)
      .then((response) => response.text())
      .then((TaskYaml) => dispatch({
        type: FETCH_TASK_YAML,
        payload: TaskYaml
      }))

  };


}

export default fetchTaskDescription;
