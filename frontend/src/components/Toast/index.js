import { Store } from 'react-notifications-component';

const Toast = ({ type, title, message }) => {
    return Store.addNotification({
        type,
        title,
        message,
        insert: "top",
        container: "top-right",
        animationIn: ["animate__animated", "animate__fadeIn"],
        animationOut: ["animate__animated", "animate__fadeOut"],
        dismiss: {
            duration: 50000,
            onScreen: true
        }
    });
}

export default Toast