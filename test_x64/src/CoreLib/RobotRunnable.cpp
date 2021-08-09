#include "RobotRunnable.h"
#include "RobotManager.h"

#include <QThread>


RobotRunnable::RobotRunnable(RobotManager* rmgr, const QString& strHost, int nPort, const QString& strPrefix, int nS, int nE)
	: QRunnable(),m_rmgr(rmgr),m_strHost(strHost),m_nPort(nPort),m_strPrefix(strPrefix),m_nS(nS),m_nE(nE)
{
	m_keylst.clear();

	int maxsz = 6;

	for (int k = nS; k < nE; ++k)
	{
		QString strSn = QString::number(k);
		int sz = maxsz - strSn.size();
		if (sz > 0)
			strSn = QString().fill('0', sz) + strSn;
		QString strUser = strPrefix + strSn;

		m_keylst.append(strUser);
	}

	//qDebug() << "Thread id = " << QThread::currentThreadId() << ", nS = " << nS << ", nE = " << nE << ", keylst.size = " << m_keylst.size();
}

RobotRunnable::~RobotRunnable()
{
}

void RobotRunnable::run()
{
	if (m_rmgr)
		m_rmgr->doRunnable(m_strHost, m_nPort, m_keylst);
}

