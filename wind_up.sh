#!/bin/bash

function main() {
    develop_branch=$1
    # 检出dev分支
    git checkout dev
    # 将需求分支合并至dev分支
    git merge  $develop_branch
    # 将合并后的dev分支打一个dev_tag
    git tag -a "dev_tag_${develop_branch}" -m "将分支${develop_branch}合并dev后，对dev打一个tag"
    # 检出master分支
    git checkout master
    # 将合并需求分支后的dev分支合并至master分支
    git merge dev
    # 将合并后的master分支打一个master_tag
    git tag -a "master_tag_${develop_branch}" -m "将分支${develop_branch}合并master后，对master打一个tag"
}

# 1、todo-done 首先确定本地已经检出了本地分支
dev_branch=`git branch | grep "dev$"`
if [ ${dev_branch} ]; then
    echo "本地已存在dev分支, 后续直接执行checkout操作"
else
    git checkout -b dev origin/dev
fi
# 2、todo-done 其次确定本地已经检出了master分支
master_branch=`git branch | grep "master$"`
if [ ${master_branch} ]; then
    echo "本地已存在master分支, 后续直接执行checkout操作"
else
    git checkout -b master origin/master
fi
# 3、todo-done 最后验证脚本执行参数
if [ $1 ]; then
    main $1
else
    echo "请输入您的需求分支，否则无法执行后续操作"
fi