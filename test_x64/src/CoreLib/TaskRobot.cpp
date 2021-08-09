#include "TaskRobot.h"

#include "RobotManager.h"


TaskRobot::TaskRobot(RobotManager* rmgr, QObject *parent)
	: TaskBase(parent), m_rmgr(rmgr)
{
	//toStartup();
}

TaskRobot::~TaskRobot()
{
}

void TaskRobot::doWork()
{
	
}


