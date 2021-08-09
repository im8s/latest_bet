#include "RobotManager.h"

#include "TaskRobot.h"
#include "ARobot.h"

#include "RobotRunnable.h"

#include <QCoreApplication>
#include <QMessageBox>
#include <QThreadPool>


RobotManager::RobotManager(QObject *parent)
	: QObject(parent)
{
	qRegisterMetaType<tPlayer>("tPlayer");
	qRegisterMetaType<tPlayer>("tPlayer&");

	qRegisterMetaType<tPlayerRefMap>("tPlayerRefMap");
	qRegisterMetaType<tPlayerRefMap>("tPlayerRefMap&");

	qRegisterMetaType<tRobot>("tRobot");
	qRegisterMetaType<tRobot>("tRobot&");

	qRegisterMetaType<tRobotVector>("tRobotVector");
	qRegisterMetaType<tRobotVector>("tRobotVector&");

	qRegisterMetaType<tRobotMap>("tRobotMap");
	qRegisterMetaType<tRobotMap>("tRobotMap&");

	/*qRegisterMetaType<tLottery>("tLottery");
	qRegisterMetaType<tLottery>("tLottery&");

	qRegisterMetaType<tLotteryRefMap>("tLotteryRefMap");
	qRegisterMetaType<tLotteryRefMap>("tLotteryRefMap&");

	qRegisterMetaType<tServerTime>("tServerTime");
	qRegisterMetaType<tServerTime>("tServerTime&");*/

	m_rbtLock = new QMutex();
	
	m_threadpool = new QThreadPool();

	qsrand(QTime(0, 0, 0).secsTo(QTime::currentTime()));
}

RobotManager::~RobotManager()
{
	m_bQuit = true;
	m_bReset = true;

	m_threadpool->waitForDone();
	delete m_threadpool;

	delete m_rbtLock;
}

void RobotManager::toDoWork(const QString& strHost, int nPort, int nThreadNum, const QString& strPrefix, int nRobotNum)
{
	m_bQuit = false;
	m_bReset = false;

	m_strHostName = strHost;
	m_nPort = nPort;
	m_nThreadNum = nThreadNum;

	m_strPrefix = strPrefix;
	m_nRobotNum = nRobotNum;

	m_threadpool->setMaxThreadCount(nThreadNum);

	int nRunNum = nThreadNum;
	int gap = nRobotNum / nRunNum;

	if (gap <= 0)
	{
		gap = 1;
		nRunNum = nRobotNum;
	}

	int nS = 0;
	int nE = nS;

	for (int k = 0; k < nRunNum; ++k)
	{
		nS = nE;
		nE = nS + gap;

		RobotRunnable* task = new RobotRunnable(this,strHost,nPort,strPrefix,nS,nE);
		m_threadpool->start(task);
	}
}

void RobotManager::doRunnable(const QString& strHost, int nPort, const QStringList& keylst)
{
	{
		try
		{
			for (int k = 0; k < keylst.size(); ++k)
			{
				{
					QString strUser = keylst[k];

					QMutexLocker locker(m_rbtLock);

					if (!m_rbtColl.contains(strUser))
					{
						ARobot* p = new ARobot();

						if (p)
						{
							p->setUserNameAndPassword(strUser, strUser);
							p->setParams(strHost, nPort);

							p->setRobotStatus(RS_init);

							m_rbtColl[strUser] = p;
						}
					}
				}

				QCoreApplication::processEvents();
			}

			getRobotColl(m_trbtColl);

			refreshModel();
		}
		catch (const bad_alloc& e)
		{
			QMessageBox::information(0, ZN_STR("信息提示"), ZN_STR("内存分配失败: %1").arg(e.what()));
		}
		catch (...)
		{
			qDebug() << "other exception";
		}
	}

	int ind = 0;

	while (!m_bQuit)
	{
		if (m_bReset)
		{
			clearAll(keylst);
			refreshModel();

			break;
		}
		else
		{
			if (ind < 0 || ind >= keylst.size())
				ind = 0;

			if (ind >= 0 && ind < keylst.size())
			{
				QString pid = keylst[ind];

				//bool b;
				{
					QMutexLocker locker(m_rbtLock);

					ARobot* p = m_rbtColl[pid];

					p->doAction();

					//b = p->resetStateChanged();
				}

				//if(b)
				//{
				//	//refreshModel();
				//}
			}

			++ind;
		}

		//QThread::msleep(1000);
		QCoreApplication::processEvents();
	}

	clearAll(keylst);

	if (!m_bQuit)
		Q_EMIT sigRobotTaskQuit();
}

void RobotManager::clearAll(const QStringList& keylst)
{
	{
		ARobotVector coll;
		{
			QMutexLocker locker(m_rbtLock);

			m_trbtColl.clear();

			for(int k = 0; k < keylst.size(); ++k)
			{
				QString pid = keylst[k];
				
				if (m_rbtColl.contains(pid))
				{
					ARobot* p = m_rbtColl[pid];
					coll.append(p);

					m_rbtColl.remove(pid);
				}
			}
		}

		for (int k = 0; k < coll.size(); ++k)
		{
			ARobot* p = coll[k];
			if (p)
			{
				p->disconnect();
				delete p;
			}

			//QCoreApplication::processEvents();
		}

		coll.clear();
	}
}


