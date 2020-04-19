import { FETCH_RESOURCE_LIST } from '../Actions/TaskActionType';
import { API_URL } from '../../../constants';

// eslint-disable-next-line require-jsdoc
export function fetchResourceList() {
    return function (dispatch: any) {
        fetch(`https://api-tekton-hub.apps.cluster-blr-65f3.blr-65f3.example.opentlc.com/resources?limit=20`)
            .then((response) => response.json())
            .then((TaskData) => dispatch({
                type: FETCH_RESOURCE_LIST,
                payload: TaskData,
            }));
    };
}

export default fetchResourceList;
;
