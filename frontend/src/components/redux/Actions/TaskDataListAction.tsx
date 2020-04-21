import { FETCH_TASK_LIST } from '../Actions/TaskActionType';
import { API_URL } from '../../../constants';

// eslint-disable-next-line require-jsdoc
export function fetchTaskList() {
  return function (dispatch: any) {
    fetch(`https://api-tekton-hub.apps.cluster-blr-8fcf.blr-8fcf.example.opentlc.com/resources?limit=5`)
      .then((response) => response.json())
      .then((TaskData) =>
        dispatch({
          type: FETCH_TASK_LIST,
          payload: TaskData.data,
        }));
  };
}

export default fetchTaskList
  ;
