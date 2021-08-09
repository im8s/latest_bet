#pragma once

#include "TaskBase.h"


class RobotManager;
class TaskRobot : public TaskBase
{
	Q_OBJECT

public:
	TaskRobot(RobotManager* rmgr, QObject *parent = nullptr);
	~TaskRobot();

protected:
	void doWork();

private:
	RobotManager*	m_rmgr = nullptr;
};
