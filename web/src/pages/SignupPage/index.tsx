import React, { Component } from 'react'
import { signup } from '../../utils/AuthService'
import MessageBar, { MessageType } from '../../components/MessageBar'
import {Link} from "react-router-dom";

export default class SignupPage extends Component<any> {
    state = {
        email: '',
        username: '',
        password: '',
        errorMsg: '',
        successMsg: '',
    };

    handleSubmit = (ev: React.SyntheticEvent<any>) => {
        ev.preventDefault();

        signup(this.state.email, this.state.username,this.state.password)
            .then((res: JSON) =>{
                this.setState({
                    successMsg: "Check your email to verify your account!",
                })
            })
            .catch((err: Error) =>{
                this.setState({
                    errorMsg: err.message,
                })
            })
    };

    render() {
        const formStyles = "bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4";
        const labelStyles = "block text-grey-darker text-sm font-bold mb-2";
        const inputStyles = "shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline bg-grey-lightest";
        const btnStyles = "bg-blue hover:bg-blue-dark text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline";

        return (
            <div className="container mx-auto m-5">
                <div className="flex justify-center">
                    <form className={formStyles} onSubmit={this.handleSubmit}>
                        {
                            this.state.successMsg !== '' ?
                                <div className="mb-4 p-4">
                                    <MessageBar title="Success!" message={this.state.successMsg} type={MessageType.Success} />
                                </div> 
                            : undefined
                        }

                        <div className="mb-4">
                            <label className={labelStyles}>
                                Email
                            </label>

                            {/* Email Input */}
                            <input 
                                type="email"
                                className={inputStyles} 
                                value={this.state.email} 
                                onChange={(ev) => this.setState({ email: ev.target.value })} />
                        </div>
                        <div className="mb-4">
                            <label className={labelStyles}>
                                Username
                            </label>

                            {/* Username Input */}
                            <input 
                                className={inputStyles} 
                                value={this.state.username} 
                                onChange={(ev) => this.setState({ username: ev.target.value })} />
                        </div>
                        <div className="mb-6">
                            <label className={labelStyles}>
                                Password
                            </label>

                            {/* Password Input */}
                            <input 
                                className={inputStyles} 
                                type="password" 
                                value={this.state.password} 
                                onChange={(ev) => this.setState({ password: ev.target.value })} />

                            {/* Error Message */}
                            { 
                                this.state.errorMsg != '' ? 
                                    <p className="text-red text-xs italic">{this.state.errorMsg}</p> 
                                    : undefined 
                            }
                        </div>
                        <div className="flex items-center justify-between">
                            <button className={btnStyles} type="submit">
                                Signup
                            </button>
                            <Link to="/login" className="inline-block align-baseline font-bold text-sm text-blue hover:text-blue-darker p-1">
                                Already have an account?
                            </Link>
                        </div>
                    </form>
                </div>
            </div>
        )
    }
}