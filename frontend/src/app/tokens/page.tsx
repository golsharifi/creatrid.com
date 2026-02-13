"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { DollarSign, Users, Activity, Send, Star, X } from "@/components/icons";
import { useTranslation } from "react-i18next";

type TokenData = {
  id: string;
  userId: string;
  name: string;
  symbol: string;
  description: string | null;
  totalSupply: number;
  priceCents: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
};

export default function TokensPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState<"token" | "holders" | "transactions" | "tips" | "fans">("token");
  const [token, setToken] = useState<TokenData | null>(null);
  const [hasToken, setHasToken] = useState<boolean | null>(null);

  // Create token form
  const [createName, setCreateName] = useState("");
  const [createSymbol, setCreateSymbol] = useState("");
  const [createDescription, setCreateDescription] = useState("");
  const [createPrice, setCreatePrice] = useState("1.00");
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState("");

  // Holders
  const [holders, setHolders] = useState<any[]>([]);
  const [holdersTotal, setHoldersTotal] = useState(0);
  const [holdersPage, setHoldersPage] = useState(0);

  // Transactions
  const [transactions, setTransactions] = useState<any[]>([]);
  const [transactionsTotal, setTransactionsTotal] = useState(0);
  const [transactionsPage, setTransactionsPage] = useState(0);

  // Tips
  const [tipsReceived, setTipsReceived] = useState<any[]>([]);
  const [tipsSent, setTipsSent] = useState<any[]>([]);
  const [tipStats, setTipStats] = useState<{ totalReceivedCents: number; totalSentCents: number; receivedCount: number; sentCount: number } | null>(null);

  // Fans
  const [fans, setFans] = useState<any[]>([]);
  const [fansTotal, setFansTotal] = useState(0);
  const [fansPage, setFansPage] = useState(0);

  // Subscriptions
  const [mySubscriptions, setMySubscriptions] = useState<any[]>([]);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  // Load token data
  useEffect(() => {
    if (!user) return;
    api.tokens.get().then((r) => {
      if (r.data && r.data.token) {
        setToken(r.data.token);
        setHasToken(true);
      } else {
        setHasToken(false);
      }
    });
  }, [user]);

  // Load holders when tab changes
  useEffect(() => {
    if (activeTab === "holders" && token) {
      api.tokens.holders(token.id, 20, holdersPage * 20).then((r) => {
        if (r.data) {
          setHolders(r.data.holders || []);
          setHoldersTotal(r.data.total);
        }
      });
    }
  }, [activeTab, token, holdersPage]);

  // Load transactions when tab changes
  useEffect(() => {
    if (activeTab === "transactions" && token) {
      api.tokens.transactions(token.id, 20, transactionsPage * 20).then((r) => {
        if (r.data) {
          setTransactions(r.data.transactions || []);
          setTransactionsTotal(r.data.total);
        }
      });
    }
  }, [activeTab, token, transactionsPage]);

  // Load tips when tab changes
  useEffect(() => {
    if (activeTab === "tips" && user) {
      Promise.all([
        api.tips.received(),
        api.tips.sent(),
        api.tips.stats(),
      ]).then(([receivedRes, sentRes, statsRes]) => {
        if (receivedRes.data) setTipsReceived(receivedRes.data.tips || []);
        if (sentRes.data) setTipsSent(sentRes.data.tips || []);
        if (statsRes.data) setTipStats(statsRes.data);
      });
    }
  }, [activeTab, user]);

  // Load fans when tab changes
  useEffect(() => {
    if (activeTab === "fans" && user) {
      api.fanSubscriptions.fans(20, fansPage * 20).then((r) => {
        if (r.data) {
          setFans(r.data.fans || []);
          setFansTotal(r.data.total);
        }
      });
      api.fanSubscriptions.list().then((r) => {
        if (r.data) {
          setMySubscriptions(r.data.subscriptions || []);
        }
      });
    }
  }, [activeTab, user, fansPage]);

  async function handleCreate() {
    if (!createName.trim() || !createSymbol.trim()) return;
    setCreating(true);
    setCreateError("");
    const priceCents = Math.round(parseFloat(createPrice) * 100) || 100;
    const result = await api.tokens.create({
      name: createName.trim(),
      symbol: createSymbol.trim().toUpperCase(),
      description: createDescription.trim() || undefined,
      priceCents,
    });
    if (result.data) {
      setToken(result.data.token);
      setHasToken(true);
    } else {
      setCreateError(result.error || "Failed to create token");
    }
    setCreating(false);
  }

  async function handleCancelSubscription(subId: string) {
    const result = await api.fanSubscriptions.cancel(subId);
    if (result.data) {
      setMySubscriptions((prev) => prev.filter((s) => s.id !== subId));
    }
  }

  if (loading || !user) return null;

  const tabs = [
    { key: "token" as const, label: t("tokens.title") },
    { key: "holders" as const, label: t("tokens.holders") },
    { key: "transactions" as const, label: t("tokens.transactions") },
    { key: "tips" as const, label: t("tokens.tips") },
    { key: "fans" as const, label: t("tokens.fans") },
  ];

  // Show create token form if no token and on token tab
  if (hasToken === false && activeTab === "token") {
    return (
      <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
        <div className="mb-8">
          <h1 className="text-2xl font-bold">{t("tokens.title")}</h1>
          <p className="mt-1 text-zinc-500">{t("tokens.subtitle")}</p>
        </div>

        {/* Tabs - still show them so user can navigate to tips/fans */}
        <div className="mb-6 flex gap-4 border-b border-zinc-200 dark:border-zinc-800">
          {tabs.map((tab) => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={`border-b-2 px-4 py-2 text-sm font-medium ${
                activeTab === tab.key
                  ? "border-zinc-900 text-zinc-900 dark:border-zinc-100 dark:text-zinc-100"
                  : "border-transparent text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300"
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        <div className="mx-auto max-w-xl">
          <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
            <h2 className="mb-4 text-lg font-semibold">{t("tokens.createToken")}</h2>
            <p className="mb-6 text-sm text-zinc-500">{t("tokens.noToken")}</p>
            <div className="space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium">{t("tokens.tokenName")}</label>
                <input
                  type="text"
                  value={createName}
                  onChange={(e) => setCreateName(e.target.value)}
                  className="w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
                  placeholder="My Creator Token"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">{t("tokens.tokenSymbol")}</label>
                <input
                  type="text"
                  value={createSymbol}
                  onChange={(e) => setCreateSymbol(e.target.value.toUpperCase())}
                  className="w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm uppercase dark:border-zinc-700 dark:bg-zinc-900"
                  placeholder="MCT"
                  maxLength={10}
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">{t("tokens.pricePerToken")} ($)</label>
                <input
                  type="number"
                  value={createPrice}
                  onChange={(e) => setCreatePrice(e.target.value)}
                  className="w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
                  placeholder="1.00"
                  min="0.01"
                  step="0.01"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">{t("tokens.description")}</label>
                <textarea
                  value={createDescription}
                  onChange={(e) => setCreateDescription(e.target.value)}
                  className="w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
                  rows={3}
                  placeholder="Describe your token..."
                />
              </div>
              {createError && (
                <div className="flex items-center gap-2 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600 dark:bg-red-900/20 dark:text-red-400">
                  <X className="h-4 w-4 shrink-0" />
                  {createError}
                </div>
              )}
              <button
                onClick={handleCreate}
                disabled={creating || !createName.trim() || !createSymbol.trim()}
                className="w-full rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
              >
                {creating ? t("common.loading") : t("tokens.createToken")}
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (hasToken === null && activeTab === "token") return null;

  return (
    <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold">{t("tokens.title")}</h1>
        <p className="mt-1 text-zinc-500">{t("tokens.subtitle")}</p>
      </div>

      {/* Tabs */}
      <div className="mb-6 flex gap-4 border-b border-zinc-200 dark:border-zinc-800">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`border-b-2 px-4 py-2 text-sm font-medium ${
              activeTab === tab.key
                ? "border-zinc-900 text-zinc-900 dark:border-zinc-100 dark:text-zinc-100"
                : "border-transparent text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300"
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Token Overview Tab */}
      {activeTab === "token" && token && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <StatCard icon={Star} label={t("tokens.tokenName")} value={`${token.name} (${token.symbol})`} />
          <StatCard icon={DollarSign} label={t("tokens.pricePerToken")} value={`$${(token.priceCents / 100).toFixed(2)}`} />
          <StatCard icon={Activity} label={t("tokens.supply")} value={token.totalSupply} />
          <StatCard icon={Users} label={t("tokens.holders")} value={holdersTotal || 0} />
        </div>
      )}

      {/* Holders Tab */}
      {activeTab === "holders" && (
        <div className="space-y-6">
          {holders.length === 0 ? (
            <div className="rounded-xl border border-zinc-200 px-6 py-12 text-center dark:border-zinc-800">
              <p className="text-zinc-500">{t("tokens.noFans")}</p>
            </div>
          ) : (
            <div className="rounded-xl border border-zinc-200 dark:border-zinc-800">
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead className="border-b border-zinc-200 text-xs text-zinc-500 dark:border-zinc-800">
                    <tr>
                      <th className="px-6 py-3 font-medium">{t("common.creator")}</th>
                      <th className="px-6 py-3 font-medium">{t("tokens.balance")}</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
                    {holders.map((h: any, i: number) => (
                      <tr key={i}>
                        <td className="px-6 py-3">{h.userName || h.userId}</td>
                        <td className="px-6 py-3">{h.balance}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              {holdersTotal > 20 && (
                <Pagination page={holdersPage} setPage={setHoldersPage} total={holdersTotal} perPage={20} />
              )}
            </div>
          )}
        </div>
      )}

      {/* Transactions Tab */}
      {activeTab === "transactions" && (
        <div className="space-y-6">
          {transactions.length === 0 ? (
            <div className="rounded-xl border border-zinc-200 px-6 py-12 text-center dark:border-zinc-800">
              <p className="text-zinc-500">{t("tokens.noToken")}</p>
            </div>
          ) : (
            <div className="rounded-xl border border-zinc-200 dark:border-zinc-800">
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead className="border-b border-zinc-200 text-xs text-zinc-500 dark:border-zinc-800">
                    <tr>
                      <th className="px-6 py-3 font-medium">{t("tokens.amount")}</th>
                      <th className="px-6 py-3 font-medium">{t("common.creator")}</th>
                      <th className="px-6 py-3 font-medium">{t("earnings.date")}</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
                    {transactions.map((tx: any, i: number) => (
                      <tr key={i}>
                        <td className="px-6 py-3 font-medium">{tx.amount}</td>
                        <td className="px-6 py-3">{tx.userName || tx.userId}</td>
                        <td className="px-6 py-3 text-zinc-500">{new Date(tx.createdAt).toLocaleDateString()}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              {transactionsTotal > 20 && (
                <Pagination page={transactionsPage} setPage={setTransactionsPage} total={transactionsTotal} perPage={20} />
              )}
            </div>
          )}
        </div>
      )}

      {/* Tips Tab */}
      {activeTab === "tips" && (
        <div className="space-y-6">
          {/* Stats row */}
          {tipStats && (
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <StatCard icon={DollarSign} label={t("tokens.totalReceived")} value={`$${(tipStats.totalReceivedCents / 100).toFixed(2)}`} />
              <StatCard icon={Send} label={t("tokens.totalSent")} value={`$${(tipStats.totalSentCents / 100).toFixed(2)}`} />
              <StatCard icon={Activity} label={t("tokens.received")} value={tipStats.receivedCount} />
              <StatCard icon={Activity} label={t("tokens.sent")} value={tipStats.sentCount} />
            </div>
          )}

          {/* Tips Received */}
          <div>
            <h3 className="mb-3 text-lg font-semibold">{t("tokens.received")}</h3>
            {tipsReceived.length === 0 ? (
              <div className="rounded-xl border border-zinc-200 px-6 py-8 text-center dark:border-zinc-800">
                <p className="text-zinc-500">{t("tokens.noTips")}</p>
              </div>
            ) : (
              <div className="rounded-xl border border-zinc-200 dark:border-zinc-800">
                <div className="overflow-x-auto">
                  <table className="w-full text-left text-sm">
                    <thead className="border-b border-zinc-200 text-xs text-zinc-500 dark:border-zinc-800">
                      <tr>
                        <th className="px-6 py-3 font-medium">{t("tokens.amount")}</th>
                        <th className="px-6 py-3 font-medium">{t("tokens.tipMessage")}</th>
                        <th className="px-6 py-3 font-medium">{t("earnings.status")}</th>
                        <th className="px-6 py-3 font-medium">{t("earnings.date")}</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
                      {tipsReceived.map((tip: any, i: number) => (
                        <tr key={i}>
                          <td className="px-6 py-3 font-medium">${(tip.amountCents / 100).toFixed(2)}</td>
                          <td className="px-6 py-3 text-zinc-500">{tip.message || t("common.noData")}</td>
                          <td className="px-6 py-3">
                            <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${
                              tip.status === "completed" ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400" : "bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400"
                            }`}>
                              {tip.status}
                            </span>
                          </td>
                          <td className="px-6 py-3 text-zinc-500">{new Date(tip.createdAt).toLocaleDateString()}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>

          {/* Tips Sent */}
          <div>
            <h3 className="mb-3 text-lg font-semibold">{t("tokens.sent")}</h3>
            {tipsSent.length === 0 ? (
              <div className="rounded-xl border border-zinc-200 px-6 py-8 text-center dark:border-zinc-800">
                <p className="text-zinc-500">{t("tokens.noTips")}</p>
              </div>
            ) : (
              <div className="rounded-xl border border-zinc-200 dark:border-zinc-800">
                <div className="overflow-x-auto">
                  <table className="w-full text-left text-sm">
                    <thead className="border-b border-zinc-200 text-xs text-zinc-500 dark:border-zinc-800">
                      <tr>
                        <th className="px-6 py-3 font-medium">{t("tokens.amount")}</th>
                        <th className="px-6 py-3 font-medium">{t("tokens.tipMessage")}</th>
                        <th className="px-6 py-3 font-medium">{t("earnings.status")}</th>
                        <th className="px-6 py-3 font-medium">{t("earnings.date")}</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
                      {tipsSent.map((tip: any, i: number) => (
                        <tr key={i}>
                          <td className="px-6 py-3 font-medium">${(tip.amountCents / 100).toFixed(2)}</td>
                          <td className="px-6 py-3 text-zinc-500">{tip.message || t("common.noData")}</td>
                          <td className="px-6 py-3">
                            <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${
                              tip.status === "completed" ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400" : "bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400"
                            }`}>
                              {tip.status}
                            </span>
                          </td>
                          <td className="px-6 py-3 text-zinc-500">{new Date(tip.createdAt).toLocaleDateString()}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Fans Tab */}
      {activeTab === "fans" && (
        <div className="space-y-6">
          {/* My Fans */}
          <div>
            <h3 className="mb-3 text-lg font-semibold">{t("tokens.fans")}</h3>
            {fans.length === 0 ? (
              <div className="rounded-xl border border-zinc-200 px-6 py-8 text-center dark:border-zinc-800">
                <p className="text-zinc-500">{t("tokens.noFans")}</p>
              </div>
            ) : (
              <div className="rounded-xl border border-zinc-200 dark:border-zinc-800">
                <div className="overflow-x-auto">
                  <table className="w-full text-left text-sm">
                    <thead className="border-b border-zinc-200 text-xs text-zinc-500 dark:border-zinc-800">
                      <tr>
                        <th className="px-6 py-3 font-medium">{t("common.creator")}</th>
                        <th className="px-6 py-3 font-medium">{t("tokens.tier")}</th>
                        <th className="px-6 py-3 font-medium">{t("earnings.status")}</th>
                        <th className="px-6 py-3 font-medium">{t("earnings.date")}</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
                      {fans.map((fan: any, i: number) => (
                        <tr key={i}>
                          <td className="px-6 py-3">{fan.fanName || fan.fanUserID || t("common.noData")}</td>
                          <td className="px-6 py-3">
                            <TierBadge tier={fan.tier} />
                          </td>
                          <td className="px-6 py-3">
                            <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${
                              fan.status === "active" ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400" : "bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400"
                            }`}>
                              {fan.status}
                            </span>
                          </td>
                          <td className="px-6 py-3 text-zinc-500">{new Date(fan.startedAt).toLocaleDateString()}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
                {fansTotal > 20 && (
                  <Pagination page={fansPage} setPage={setFansPage} total={fansTotal} perPage={20} />
                )}
              </div>
            )}
          </div>

          {/* My Subscriptions */}
          <div>
            <h3 className="mb-3 text-lg font-semibold">{t("tokens.mySubscriptions")}</h3>
            {mySubscriptions.length === 0 ? (
              <div className="rounded-xl border border-zinc-200 px-6 py-8 text-center dark:border-zinc-800">
                <p className="text-zinc-500">{t("tokens.noSubscriptions")}</p>
              </div>
            ) : (
              <div className="rounded-xl border border-zinc-200 dark:border-zinc-800">
                <div className="overflow-x-auto">
                  <table className="w-full text-left text-sm">
                    <thead className="border-b border-zinc-200 text-xs text-zinc-500 dark:border-zinc-800">
                      <tr>
                        <th className="px-6 py-3 font-medium">{t("common.creator")}</th>
                        <th className="px-6 py-3 font-medium">{t("tokens.tier")}</th>
                        <th className="px-6 py-3 font-medium">{t("earnings.status")}</th>
                        <th className="px-6 py-3 font-medium">{t("admin.tableActions")}</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
                      {mySubscriptions.map((sub: any, i: number) => (
                        <tr key={i}>
                          <td className="px-6 py-3">{sub.creatorName || sub.creatorUserID || t("common.noData")}</td>
                          <td className="px-6 py-3">
                            <TierBadge tier={sub.tier} />
                          </td>
                          <td className="px-6 py-3">
                            <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${
                              sub.status === "active" ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400" : "bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400"
                            }`}>
                              {sub.status}
                            </span>
                          </td>
                          <td className="px-6 py-3">
                            {sub.status === "active" && (
                              <button
                                onClick={() => handleCancelSubscription(sub.id)}
                                className="rounded-lg border border-red-200 px-3 py-1 text-xs font-medium text-red-600 transition-colors hover:bg-red-50 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-900/30"
                              >
                                {t("tokens.cancelSubscription")}
                              </button>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

function StatCard({
  icon: Icon,
  label,
  value,
}: {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  value: number | string;
}) {
  return (
    <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
      <div className="mb-2 flex h-9 w-9 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
        <Icon className="h-4 w-4" />
      </div>
      <p className="text-sm text-zinc-500">{label}</p>
      <p className="mt-1 text-2xl font-bold">{typeof value === "number" ? value.toLocaleString() : value}</p>
    </div>
  );
}

function TierBadge({ tier }: { tier: string }) {
  const colors: Record<string, string> = {
    supporter: "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400",
    superfan: "bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400",
    patron: "bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400",
  };
  return (
    <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium capitalize ${colors[tier] || "bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400"}`}>
      {tier}
    </span>
  );
}

function Pagination({
  page,
  setPage,
  total,
  perPage,
}: {
  page: number;
  setPage: (fn: (p: number) => number) => void;
  total: number;
  perPage: number;
}) {
  const { t } = useTranslation();
  return (
    <div className="flex items-center justify-between border-t border-zinc-200 px-6 py-3 dark:border-zinc-800">
      <button
        onClick={() => setPage((p) => Math.max(0, p - 1))}
        disabled={page === 0}
        className="rounded-lg border border-zinc-200 px-3 py-1 text-xs font-medium disabled:opacity-50 dark:border-zinc-700"
      >
        {t("common.previous")}
      </button>
      <span className="text-xs text-zinc-500">
        {t("common.page", { current: page + 1, total: Math.ceil(total / perPage) })}
      </span>
      <button
        onClick={() => setPage((p) => p + 1)}
        disabled={(page + 1) * perPage >= total}
        className="rounded-lg border border-zinc-200 px-3 py-1 text-xs font-medium disabled:opacity-50 dark:border-zinc-700"
      >
        {t("common.next")}
      </button>
    </div>
  );
}
