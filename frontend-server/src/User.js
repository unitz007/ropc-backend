import './App.css';

class User {
    constructor(username, email, clientId) {
        this.username = username;
        this.email = email;
        this.clientId = clientId;
    }
}

// function  FetchUserDetails() {
//
//     fetch()
//
//     return new User()
// }

function UserDetails() {

    const user = new User("charles", "dinneyacharles007gmail.com", "clientId")

    return (
        <div className="container">
            <h1 className="btn-secondary">Profile</h1>
            &nbsp;
            <table className="table">
                <tbody>
                <tr>
                    <th>Username</th>
                    <td>{user.username}</td>
                </tr>
                <tr>
                    <th>Client Id</th>
                    <td>{user.clientId}</td>
                </tr>
                <tr>
                    <th>Email Address</th>
                    <td>{user.email}</td>
                </tr>
                </tbody>
            </table>
            <h1 className="btn-secondary">&nbsp;</h1>
        </div>
    )
}

export default UserDetails;