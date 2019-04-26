import React, { Component } from 'react'
import { RouteComponentProps } from 'react-router'
import { ClipLoader } from 'react-spinners'

import { verifyEmail } from '../../utils/AuthService'

interface MatchParams {
    token: string;
}

interface Props extends RouteComponentProps<MatchParams> {}

export default class VerifyEmailPage extends Component<Props> {

    state = {
        isLoading: true,
        isSuccessful: false,
    }

    constructor(props: Props) {
        super(props)
    }

    componentDidMount() {
        // Get token
        let token = this.props.match.params.token;

        // Make API request
        let res = verifyEmail(token);
        
        // When ready, set state to given response
        res.then(r => {
            let isSuccess = r.status === 200;

            this.setState({
                isLoading: false,
                isSuccessful: isSuccess,
            })
        })
    }

    render() {
        return (
            <div className="text-center">
                <ClipLoader
                    sizeUnit={"px"}
                    size={150}
                    color={'#123abc'}
                    loading={this.state.isLoading} />
                    
                <p>
                    { this.state.isLoading ? undefined : (this.state.isSuccessful ? 'Success' : 'Failure') }
                </p>
                <i>TODO: Show a 'Resend email verification' link in case of failure.</i>
            </div>
        )
    }
}
