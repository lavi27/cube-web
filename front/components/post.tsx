interface compType {
    userIconSrc: string;
    userName: string;
    date: number;
    content: string;
    like: number;
}

const comp = ({userIconSrc, userName, date, content, like}: compType) => {
    return (
        <div className="post">
            <div className="post_header">
                <div className="user_wrap">
                    <div className="user_icon_wrap">
                        <img src={userIconSrc}></img>
                    </div>
                    <span className="user_userName">{userName}</span>
                </div>
                <span className="date">{date}</span>
            </div>
            <div className="post_content">{content}</div>
            <div className="post_footer">
                <div className="like_wrap">
                    <div className="likeImg_wrap">
                        <img src="heart"></img>
                    </div>
                    <span>{like}</span>
                </div>
            </div>
        </div>
    )
}

export default comp;