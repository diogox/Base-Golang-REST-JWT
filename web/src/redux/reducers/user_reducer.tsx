import { Action } from "../actions";
import { LOGIN } from "../constants";
import { initialState } from "./index";

export function userReducer(state = initialState, action: Action) {
    switch (action.type) {
        case LOGIN: {
            return {
                ...state,
                user: action.payload,
            }
        }
    } /* else if (action.type === LOGIN) {
        return Object.assign({}, state, {
            articles: state.articles.concat(action.payload)
        });
    }*/

    return state;
}