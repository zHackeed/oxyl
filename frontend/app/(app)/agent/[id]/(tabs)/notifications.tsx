import GlobalHeader from "@/components/ui/Header"
import { WrappedView } from "@/components/ui/WrappedView"
import { useQuery } from "@tanstack/react-query"

export default function Notifications() {

  const { data , isLoading } = useQuery({
    queryKey: ["agent-notifications"],
    queryFn: () => {
      return []
    },
  })


  return (
    <WrappedView>
      <GlobalHeader title="Notificaciones" description={`${data?.length || 0} notificaciones`} />
    </WrappedView>
  )
}