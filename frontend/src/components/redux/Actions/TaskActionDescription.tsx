import {FETCH_TASK_DESCRIPTION} from '../Actions/TaskActionType';
import {FETCH_TASK_YAML} from '../Actions/TaskActionType';
import {API_URL} from '../../../constants';

// eslint-disable-next-line require-jsdoc
export function fetchTaskDescription(id: any) {
  return function(dispatch: any) {
    fetch(`${API_URL}/resource/${id}`)
        .then((response) => response.json())
        .then(() => {
          fetch(`${API_URL}/resource/readme/${id}`)
              .then((response) => response.text())
              .then((TaskDescription) => dispatch({
                type: FETCH_TASK_DESCRIPTION,
                payload: TaskDescription,
              }));
        });

    fetch(`${API_URL}/resource/${id}`)
        .then((response) => response.json())
        .then(() => {
          fetch(`${API_URL}/resource/yaml/${id}`)
              .then((response) => response.text())
              .then((TaskYaml) => dispatch({
                type: FETCH_TASK_YAML,
                payload: TaskYaml,
              }));
        });
  };
}

export default fetchTaskDescription;
