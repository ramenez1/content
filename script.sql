create table t_content_details_3
(
    id              int auto_increment
        primary key,
    title           varchar(255) default ''                null comment '内容标题',
    description     text                                   null comment '内容描述',
    author          varchar(255) default ''                null,
    video_url       varchar(255) default ''                null,
    thumbnail       varchar(255) default ''                null comment '封面图URL',
    category        varchar(255) default ''                null comment '内容分类',
    duration        bigint       default 0                 null comment '内容时长',
    resolution      varchar(20)  default ''                null comment '分辨率 如720p、1080p',
    fileSize        bigint       default 0                 null comment '文件大小',
    format          varchar(20)  default ''                null comment '文件格式 如MP4、AVI',
    quality         int          default 0                 null comment '视频质量 1-高清 2-标清',
    approval_status int          default 0                 null comment '审核状态 1-审核中 2-审核通过 3-审核不通过',
    created_at      timestamp    default CURRENT_TIMESTAMP null comment '内容创建时间',
    updated_at      timestamp    default CURRENT_TIMESTAMP null comment '内容更新时间'
);


